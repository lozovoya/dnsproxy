package cache

import (
	"context"
	"dnproxier/server/cache/inmem"
)

type DNSCacheInterface interface {
	AddToCache(ctx context.Context, url string, response []byte) error
	GetFromCache(ctx context.Context, url string) (response []byte, err error)
	ListAllRecords(ctx context.Context) (list []string, err error)
	DeleteFromCache(ctx context.Context, url string) error
}

type DNSCache struct {
	Cache *inmem.InMemCache
}

func NewDNSCache() DNSCacheInterface {
	return &DNSCache{Cache: inmem.New()}
}

func (d *DNSCache) AddToCache(ctx context.Context, url string, response []byte) error {
	return d.Cache.AddToCache(ctx, url, response)
}

func (d *DNSCache) GetFromCache(ctx context.Context, url string) (response []byte, err error) {
	response = make([]byte, 0)
	response, err = d.Cache.GetFromCache(ctx, url)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (d *DNSCache) ListAllRecords(ctx context.Context) (list []string, err error) {
	list = make([]string, 0)
	list, err = d.Cache.ListAllRecords(ctx)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (d *DNSCache) DeleteFromCache(ctx context.Context, url string) error {
	return d.Cache.DeleteFromCache(ctx, url)
}
