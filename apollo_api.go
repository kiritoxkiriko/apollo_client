package apollo_client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type NameSpaceInfoResp struct {
	AppId         string `json:"appId"`
	ClusterName   string `json:"clusterName"`
	NamespaceName string `json:"namespaceName"`
	Comment       string `json:"comment"`
	Format        string `json:"format"`
	IsPublic      bool   `json:"isPublic"`
	Items         []struct {
		Key                        string `json:"key"`
		Value                      string `json:"value"`
		DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
		DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
		DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
		DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
	} `json:"items"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type GetConfResp struct {
	Key                        string `json:"key"`
	Value                      string `json:"value"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type AddConfReq struct {
	Key                 string `json:"key"`
	Value               string `json:"value"`
	Comment             string `json:"comment"`
	DataChangeCreatedBy string `json:"dataChangeCreatedBy"`
}

type UpdateConfReq struct {
	Key                      string `json:"key"`
	Value                    string `json:"value"`
	Comment                  string `json:"comment"`
	DataChangeLastModifiedBy string `json:"dataChangeLastModifiedBy"`
}

type AddConfResp struct {
	Key                        string `json:"key"`
	Value                      string `json:"value"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type ReleaseConfReq struct {
	ReleaseTitle   string `json:"releaseTitle"`
	ReleaseComment string `json:"releaseComment"`
	ReleasedBy     string `json:"releasedBy"`
}

type ReleaseConfResp struct {
	AppId                      string                 `json:"appId"`
	ClusterName                string                 `json:"clusterName"`
	NamespaceName              string                 `json:"namespaceName"`
	Name                       string                 `json:"name"`
	Configurations             map[string]interface{} `json:"configurations"`
	Comment                    string                 `json:"comment"`
	DataChangeCreatedBy        string                 `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string                 `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string                 `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string                 `json:"dataChangeLastModifiedTime"`
}

type GetReleasedConfResp struct {
	AppId                      string                 `json:"appId"`
	ClusterName                string                 `json:"clusterName"`
	NamespaceName              string                 `json:"namespaceName"`
	Name                       string                 `json:"name"`
	Configurations             map[string]interface{} `json:"configurations"`
	Comment                    string                 `json:"comment"`
	DataChangeCreatedBy        string                 `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string                 `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string                 `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string                 `json:"dataChangeLastModifiedTime"`
}

type ApolloApiClient struct {
	resty   *resty.Client
	creator string
}

func NewApolloApiClient(resty *resty.Client, creator string) *ApolloApiClient {
	return &ApolloApiClient{
		resty:   resty,
		creator: creator,
	}
}

func (c *ApolloApiClient) GetNameSpaceInfo(env, appId, clusterName, namespace string) (*NameSpaceInfoResp, error) {
	var resp NameSpaceInfoResp
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", env, appId, clusterName, namespace)
	rResp, err := c.resty.R().Get(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return nil, fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	err = json.Unmarshal(rResp.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ApolloApiClient) GetConf(env, appId, clusterName, namespace, key string) (*GetConfResp, error) {
	var resp GetConfResp
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s", env, appId, clusterName, namespace, key)
	rResp, err := c.resty.R().Get(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return nil, fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	err = json.Unmarshal(rResp.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ApolloApiClient) AddConf(env, appId, clusterName, namespace, key, value string) (*AddConfResp, error) {
	var resp AddConfResp
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items", env, appId, clusterName, namespace)
	req := AddConfReq{
		Key:                 key,
		Value:               value,
		Comment:             "add key: " + key,
		DataChangeCreatedBy: c.creator,
	}
	rResp, err := c.resty.R().SetBody(req).Post(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return nil, fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	err = json.Unmarshal(rResp.Body(), &resp)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	return &resp, nil
}

func (c *ApolloApiClient) UpdateConf(env, appId, clusterName, namespace, key, value string) error {
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s", env, appId, clusterName, namespace, key)
	req := UpdateConfReq{
		Key:                      key,
		Value:                    value,
		Comment:                  "update key: " + key,
		DataChangeLastModifiedBy: c.creator,
	}
	rResp, err := c.resty.R().SetBody(req).Put(path)
	if err != nil {
		return fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		fmt.Println(rResp.Request.Body)
		return fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	return nil
}

func (c *ApolloApiClient) DeleteConf(env, appId, clusterName, namespace, key string) error {
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s?operator=%s", env, appId, clusterName, namespace, key, c.creator)
	rResp, err := c.resty.R().Delete(path)
	if err != nil {
		return fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	return nil
}

func (c *ApolloApiClient) ReleaseConf(env, appId, clusterName, namespace string) (*ReleaseConfResp, error) {
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases", env, appId, clusterName, namespace)
	req := ReleaseConfReq{
		ReleaseTitle:   fmt.Sprintf("release-%d", time.Now().Unix()),
		ReleaseComment: "release conf at: " + time.Now().Format("2006-01-02 15:04:05"),
		ReleasedBy:     c.creator,
	}
	rResp, err := c.resty.R().SetBody(req).Post(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return nil, fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	var resp ReleaseConfResp
	err = json.Unmarshal(rResp.Body(), &resp)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	return &resp, nil
}

func (c *ApolloApiClient) GetReleasedConf(env, appId, clusterName, namespace string) (*GetReleasedConfResp, error) {
	path := fmt.Sprintf("/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases/latest", env, appId, clusterName, namespace)
	rResp, err := c.resty.R().Get(path)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	if rResp.StatusCode() != 200 {
		return nil, fmt.Errorf("fail to get resp: %s", rResp.String())
	}
	var resp GetReleasedConfResp
	err = json.Unmarshal(rResp.Body(), &resp)
	if err != nil {
		return nil, fmt.Errorf("fail to get resp: %s", err)
	}
	return &resp, nil
}
