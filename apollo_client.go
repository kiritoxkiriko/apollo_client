package apollo_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	logger *zap.Logger
)

type ApolloConfig struct {
	Host       string //访问地址
	AppId      string //应用id
	NameSpace  string //命名空间
	Env        string //环境信息
	Cluster    string //集群名称
	Token      string //Token
	RetryCount uint   //重试次数
	Timeout    uint   //超时时间
}

// ApolloClient an apollo client sync local to remote
type ApolloClient struct {
	Conf         ApolloConfig
	restyClient  *resty.Client
	apiClient    *ApolloApiClient
	cache        *ApolloCache
	autoSyncLock sync.Locker
	// use 3rd party lock, with try-lock function, to avoid blocking other goroutine
	// it only a temporary solution, we will use go-lock package after go 1.18 version
	// syncLock   lock.Mutex
	// use blocking lock instead of try-lock
	syncLock   sync.Locker
	isAutoSync bool
	done       chan struct{}
}

// NewApolloClient create a apollo client
func NewApolloClient(conf ApolloConfig, creator string) (*ApolloClient, error) {
	client := &ApolloClient{
		Conf:         ApolloConfig{},
		restyClient:  resty.New(),
		cache:        NewApolloCache(),
		autoSyncLock: &sync.Mutex{},
		// use 3rd party lock, with try-lock function, to avoid blocking other goroutine
		// usage: https://pkg.go.dev/github.com/viney-shih/go-lock#section-documentation
		// syncLock: lock.NewChanMutex(),
		syncLock: &sync.Mutex{},
	}
	client.SetConf(conf)
	apiClient := NewApolloApiClient(client.restyClient, creator)
	client.apiClient = apiClient

	// sync at first time
	err := client.SyncFromRemote()
	if err != nil {
		return nil, fmt.Errorf("sync from remote failed, err: %v", err)
	}
	return client, nil
}

// SetConf set apollo config
func (c *ApolloClient) SetConf(conf ApolloConfig) {
	c.Conf = conf
	// set token header
	c.restyClient.SetHeader("Authorization", c.Conf.Token)

	if c.Conf.Timeout > 0 {
		c.restyClient.SetRetryCount(int(conf.RetryCount))
	}
	if conf.Timeout > 0 {
		c.restyClient.SetTimeout(time.Duration(conf.Timeout) * time.Second)
	}
	c.restyClient.SetBaseURL(conf.Host)
}

// SetKey set key value to cache
func (c *ApolloClient) SetKey(key, value string) {
	c.cache.Set(key, value)
}

// GetKey get key value from cache
func (c *ApolloClient) GetKey(key string) (string, bool) {
	return c.cache.Get(key)
}

// UpdateKey update key value to cache
func (c *ApolloClient) UpdateKey(key, value string) {
	c.cache.Set(key, value)
}

// DeleteKey delete key from cache
func (c *ApolloClient) DeleteKey(key string) {
	c.cache.Delete(key)
}

// GetCache GetKeys get local cache
func (c *ApolloClient) GetCache() *ApolloCache {
	return c.cache
}

// SyncFromRemote this will sync the cache from remote
func (c *ApolloClient) SyncFromRemote() error {
	// if you can't get lock, return error
	//if !c.syncLock.TryLock() {
	//	return fmt.Errorf("can't get lock, another sync is running")
	//}
	//defer c.syncLock.Unlock()

	// use blocking lock instead of try-lock
	c.syncLock.Lock()
	defer c.syncLock.Unlock()

	api := c.apiClient
	info, err := api.GetReleasedConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace)
	if err != nil {
		return fmt.Errorf("get namespace info error: %v", err)
	}
	for k, v := range info.Configurations {
		c.cache.Set(k, fmt.Sprintf("%v", v))
	}
	return nil
}

// Sync this will sync local cache to remote
func (c *ApolloClient) Sync() error {
	// if you can't get lock, return error
	//if !c.syncLock.TryLock() {
	//	return fmt.Errorf("can't get lock, another sync is running")
	//}
	//defer c.syncLock.Unlock()

	// use blocking lock instead of try-lock
	c.syncLock.Lock()
	defer c.syncLock.Unlock()

	remoteValues := make(map[string]string)
	info, err := c.apiClient.GetReleasedConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace)
	if err != nil {
		return fmt.Errorf("get namespace info error: %v", err)
	}
	for k, v := range info.Configurations {
		remoteValues[k] = fmt.Sprintf("%v", v)
	}

	api := c.apiClient
	// operations count, if count is 0, means no need to update
	opsCount := 0
	var returnErr error
	c.cache.Range(func(key, value string) bool {
		// if remote don't have it , then add it
		if _, ok := remoteValues[key]; !ok {
			_, err := api.AddConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace, key, value)
			if err != nil {
				returnErr = fmt.Errorf("add key %s error: %v", key, err)
				return false
			}
			opsCount++
		} else {
			// if remote have it and local have it, then update it
			if remoteValues[key] != value {
				err := api.UpdateConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace, key, value)
				if err != nil {
					returnErr = fmt.Errorf("update key %s error: %v", key, err)
					return false
				}
			}
			opsCount++
		}
		return true
	})

	if returnErr != nil {
		return fmt.Errorf("sync error: %v", returnErr)
	}

	for key := range remoteValues {
		// if local don't have it , then delete it
		if _, ok := c.cache.Get(key); !ok {
			err := api.DeleteConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace, key)
			if err != nil {
				return fmt.Errorf("sync error: delete key %s error: %v", key, err)
			}
			opsCount++
		}
	}

	if opsCount > 0 {
		// need update, release conf
		_, err := api.ReleaseConf(c.Conf.Env, c.Conf.AppId, c.Conf.Cluster, c.Conf.NameSpace)
		if err != nil {
			return fmt.Errorf("release namespace error: %v", err)
		}
	}
	return nil
}

// AutoSync this will start auto sync, call AutoSync function every interval
func (c *ApolloClient) AutoSync(duration time.Duration) error {
	c.autoSyncLock.Lock()
	defer c.autoSyncLock.Unlock()
	if c.isAutoSync {
		return fmt.Errorf("auto sync is already running")
	}
	c.done = make(chan struct{}, 1)
	c.isAutoSync = true

	// start a go routine to sync
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-time.After(duration):
				logger.Debug("start sync")
				err := c.Sync()
				if err != nil {
					logger.Sugar().Infof("auto sync error: %v", err)
				}
				logger.Debug("sync success")
				// use sleep to avoid too frequent sync
				time.Sleep(duration)
			}
		}
	}()
	return nil
}

// StopAutoSync this will stop auto sync
func (c *ApolloClient) StopAutoSync() {
	c.autoSyncLock.Lock()
	defer c.autoSyncLock.Unlock()
	if c.isAutoSync {
		close(c.done)
		c.isAutoSync = false
	}
}

func init() {
	initLogger()
}

func initLogger() {
	logger, _ = zap.NewProduction()
}
