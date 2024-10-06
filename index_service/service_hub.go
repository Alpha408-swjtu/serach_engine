package index_service

import (
	"log"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
)

// 用etcd实现服务注册与发现
type ServiceHub struct {
	Client             *etcdv3.Client
	heartbeatFrequency int
}

var (
	serviceHub *ServiceHub
	hubOnce    sync.Once
)

func GetServiceHub(etcdServers []string, heartbeayFrequency int) *ServiceHub {
	if serviceHub == nil {
		hubOnce.Do(func() {
			if client, err := etcdv3.New(
				etcdv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				},
			); err != nil {
				log.Fatalf("连不上etcd:%s", err)
			} else {
				serviceHub = &ServiceHub{
					Client:             client,
					heartbeatFrequency: heartbeayFrequency,
				}
			}
		})
	}
	return serviceHub
}
