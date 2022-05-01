package apollo_client

import (
	"sync"
)

type ApolloCache struct {
	cache sync.Map
}

func NewApolloCache() *ApolloCache {
	return &ApolloCache{}
}

func (c *ApolloCache) Get(key string) (value string, ok bool) {
	v, ok := c.cache.Load(key)
	if !ok {
		return "", ok
	}
	return v.(string), ok
}

func (c *ApolloCache) Set(key string, value string) {
	c.cache.Store(key, value)
}

func (c *ApolloCache) Delete(key string) {
	c.cache.Delete(key)
}

func (c *ApolloCache) Range(rangeFunc func(key, value string) bool) {
	rangeFunc1 := func(key, value interface{}) bool {
		return rangeFunc(key.(string), value.(string))
	}
	c.cache.Range(rangeFunc1)
}
