package main

import (
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type Backend struct {
	URL          string
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.Alive = alive
	b.mu.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

// Ping a backend
func healthCheck(url string) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url + "/health")
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}
