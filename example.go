package apollo_client

import (
	"fmt"
	"time"
)

func main() {
	// recover from panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
		}
	}()
	// ...
	conf := ApolloConfig{
		Host:       "your host",
		AppId:      "your appid",
		NameSpace:  "your namespace",
		Env:        "DEV",
		Cluster:    "default",
		Token:      "your token",
		RetryCount: 1,
		Timeout:    2,
	}
	// create client
	// it will sync config from apollo when the client is created
	client, err := NewApolloClient(conf, "test")
	if err != nil {
		panic(err)
	}

	// set key value to local cache
	client.SetKey("key", "value")
	value, ok := client.GetKey("key")
	if ok {
		// if key is exist, it wil update the value
		client.SetKey("key", value+"1")
	}

	// sync config to apollo, return err if failed
	err = client.Sync()
	if err != nil {
		panic(err)
	}

	// get all key and value from apollo to local cache, return err if failed
	err = client.SyncFromRemote()
	if err != nil {
		panic(err)
	}

	// start auto sync, will sync local config to apollo every interval seconds
	err = client.AutoSync(10 * time.Second)
	if err != nil {
		panic(err)
	}
}
