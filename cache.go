package cache

import (
	"errors"
	"sync"
	"time"
)

type Item[T any] struct {
	Data              T
	ExpireAtTimestamp int64
}

type cache[T any] struct {
	wg   sync.WaitGroup
	mu   sync.RWMutex
	stop chan any
	data map[int64]Item[T]
}

func NewCache[T any](cleanupInterval time.Duration) *cache[T] {
	c := &cache[T]{
		wg:   sync.WaitGroup{},
		mu:   sync.RWMutex{},
		stop: make(chan any),
		data: make(map[int64]Item[T]),
	}
	c.wg.Add(1)
	go func(interval time.Duration) {
		defer c.wg.Done()
		c.cleanupLoop(interval)
	}(cleanupInterval)

	return c
}

func (c *cache[T]) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.mu.Lock()
			for uid, item := range c.data {
				if item.ExpireAtTimestamp <= time.Now().Unix() {
					delete(c.data, uid)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *cache[T]) Add(data T, id int64, expireAtTimestamp int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = Item[T]{
		Data:              data,
		ExpireAtTimestamp: expireAtTimestamp,
	}
}

func (c *cache[T]) ReadById(id int64) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var empty T
	data, ok := c.data[id]
	if !ok {
		return empty, errors.New("data not in the cache")
	}
	return data.Data, nil
}

func (c *cache[T]) ReadAll() ([]T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var result []T
	for _, data := range c.data {
		result = append(result, data.Data)
	}
	return result, nil
}

func (c *cache[T]) DeleteById(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, id)
}
