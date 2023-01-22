package inmem

import (
	"context"
	"fmt"
	"sync"
)

type InMemCache struct {
	records map[string][]byte
	mu      sync.RWMutex
}

func New() *InMemCache {
	return &InMemCache{records: make(map[string][]byte), mu: sync.RWMutex{}}
}

func (i *InMemCache) AddToCache(ctx context.Context, url string, response []byte) error {
	if len(response) == 0 {
		return fmt.Errorf("empty data")
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.records[url] = response
	return nil
}

func (i *InMemCache) GetFromCache(ctx context.Context, url string) (response []byte, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	response, ok := i.records[url]
	if !ok {
		return nil, fmt.Errorf("No record %s in cache", url)
	}
	return response, nil
}

func (i *InMemCache) ListAllRecords(ctx context.Context) (list []string, err error) {
	list = make([]string, 0) //todo так и не разобрался что там с мьютексом
	i.mu.Lock()
	defer i.mu.Unlock()
	for k, _ := range i.records {
		list = append(list, k)
	}
	return list, nil
}

func (i *InMemCache) DeleteFromCache(ctx context.Context, url string) error {
	i.mu.Lock()
	delete(i.records, url)
	i.mu.Unlock()
	return nil
}
