package sync

import (
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	Count int
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Count++
}

func (c *Counter) Value() int {
	return c.Count
}
