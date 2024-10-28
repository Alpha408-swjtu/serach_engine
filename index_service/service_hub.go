package index_service

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
)

const SERVICE_ROOT_PATH = "/alpha/index"

// 用etcd实现服务注册与发现
type ServiceHub struct {
	Client             *etcdv3.Client
	heartbeatFrequency int
	loadBalance        LoadBalance // 选择哪种负载均衡策略，在构造函数中赋值
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
					loadBalance:        &RoundRobin{}, //此处选取轮询策略
				}
			}
		})
	}
	return serviceHub
}

// 注册功能
func (hub *ServiceHub) Regist(service string, endpoint string, leaseId *etcdv3.LeaseID) (etcdv3.LeaseID, error) {
	ctx := context.Background()
	if int(*leaseId) <= 0 {
		if lease, err := hub.Client.Grant(ctx, int64(hub.heartbeatFrequency)); err != nil {
			log.Printf("创建租约失败:%v", err)
			return 0, err
		} else {
			key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
			if _, err := hub.Client.Put(ctx, key, "", etcdv3.WithLease(lease.ID)); err != nil {
				log.Printf("服务注册失败:%s", err)
				return lease.ID, err
			} else {
				return lease.ID, nil
			}
		}
	} else {
		if _, err := hub.Client.KeepAliveOnce(ctx, *leaseId); err != nil {
			return 0, err
		} else {
			return *leaseId, nil
		}
	}
}

// 注销
func (hub *ServiceHub) UnRegist(service string, endPoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endPoint
	if _, err := hub.Client.Delete(ctx, key); err != nil {
		log.Printf("注销服务:%s对应节点:%s失败:%v", service, endPoint, err)
		return err
	} else {
		log.Printf("注销成功")
		return nil
	}
}

// 返回多个服务器ip
func (hub *ServiceHub) GetServiceEndpoints(service string) []string {
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	if rep, err := hub.Client.Get(ctx, prefix, etcdv3.WithPrefix()); err != nil {
		log.Printf("获取节点失败")
		return nil
	} else {
		endpoints := make([]string, 0, len(rep.Kvs))
		for _, kv := range rep.Kvs {
			path := strings.Split(string(kv.Key), "/")
			endpoints = append(endpoints, path[len(path)-1])
		}
		log.Printf("获取节点成功")
		return endpoints
	}
}

// 经负载均衡算法,只返回一个服务器ip
func (hub *ServiceHub) GetServiceEndpoint(service string) string {
	return hub.loadBalance.Take(hub.GetServiceEndpoints(service))
}
