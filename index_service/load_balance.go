package index_service

import (
	"log"
	"math/rand"
	"sync/atomic"
)

type LoadBalance interface {
	Take([]string) string
}

// 1.轮询法
type RoundRobin struct {
	acc int64
}

func (r *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		log.Panicf("传入的服务器数量为0")
		return ""
	}
	n := atomic.AddInt64(&r.acc, 1)
	index := int64(n % int64(len(endpoints)))
	return endpoints[index]
}

// 2.随机法
type RandomSelet struct {
}

func (r *RandomSelet) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	index := rand.Intn(len(endpoints))
	return endpoints[index]
}
