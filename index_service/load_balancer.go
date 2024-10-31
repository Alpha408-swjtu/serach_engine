package index_service

import (
	"math/rand"
	"sync/atomic"
)

// 负载均衡
type LoadBalancer interface {
	Take([]string) string
}

// 1.轮询法
type RoundRobin struct {
	acc int64
}

func (b *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddInt64(&b.acc, 1)
	index := int(n % int64(len(endpoints)))
	return endpoints[index]
}

type RandomSelect struct {
}

func (r *RandomSelect) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	index := rand.Intn(len(endpoints))
	return endpoints[index]
}
