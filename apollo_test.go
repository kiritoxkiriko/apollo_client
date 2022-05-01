package apollo_client

import (
	"fmt"
	"testing"
	"time"
)

func TestMyApolloClient(t *testing.T) {
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
	client, err := NewApolloClient(conf, "creator name")

	if err != nil {
		t.Error(err)
		return // return here to avoid panic
	}
	cache := client.GetCache()
	cache.Range(func(key, value string) bool {
		t.Logf("%s = %s", key, value)
		return true
	})
	err = client.AutoSync(3 * time.Second)
	if err != nil {
		t.Error(err)
		return // return here to avoid panic
	}
	client.SetKey("test1", time.Now().String())
	client.SetKey("test2", time.Now().String())
	client.SetKey("test3", time.Now().String())
	<-time.After(5 * time.Second)
	fmt.Println("new")
	cache.Range(func(key, value string) bool {
		t.Logf("%s = %s", key, value)
		return true
	})
	client.DeleteKey("test3")
	<-time.After(5 * time.Second)
	fmt.Println("new2")
	cache.Range(func(key, value string) bool {
		t.Logf("%s = %s", key, value)
		return true
	})

}
