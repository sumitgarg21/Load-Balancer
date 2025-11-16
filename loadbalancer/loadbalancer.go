package main

import (
	"hash/fnv"
	//"math/rand"
	"sync/atomic"
)

type LoadBalancer struct {
	Backends []*Backend
	Counter  uint64
}

// Round Robin
func (lb *LoadBalancer) NextBackend() *Backend {
	next := atomic.AddUint64(&lb.Counter, 1)
	idx := int(next % uint64(len(lb.Backends)))

	// skip dead backends
	for !lb.Backends[idx].IsAlive() {
		next = atomic.AddUint64(&lb.Counter, 1)
		idx = int(next % uint64(len(lb.Backends)))
	}
	return lb.Backends[idx]
}

// Sticky Session (based on client IP hashing)
func (lb *LoadBalancer) GetBackendSticky(ip string) *Backend {
	hash := fnv.New32a()
	hash.Write([]byte(ip))
	idx := int(hash.Sum32()) % len(lb.Backends)

	if lb.Backends[idx].IsAlive() {
		return lb.Backends[idx]
	}

	// fallback to RR
	return lb.NextBackend()
}
