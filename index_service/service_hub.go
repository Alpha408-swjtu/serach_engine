package index_service

import (
	"context"
	"search_engine/utils"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

const (
	SERVICE_ROOT_PATH = "/alpha/index"
)

// 服务注册中心
type ServiceHub struct {
	client             *etcdv3.Client
	heartbeatFrequency int64 //server每隔几秒钟不动向中心上报一次心跳（其实就是续一次租约）
	watched            sync.Map
	loadBalancer       LoadBalancer //策略模式。完成同一个任务可以有多种不同的实现方案
}

var (
	serviceHub *ServiceHub //该全局变量包外不可见，包外想使用时通过GetServiceHub()获得
	hubOnce    sync.Once   //单例模式需要用到一个once
)

// ServiceHub的构造函数，单例模式
func GetServiceHub(etcdServers []string, heartbeatFrequency int64) *ServiceHub {
	if serviceHub == nil {
		hubOnce.Do(func() {
			if client, err := etcdv3.New(
				etcdv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				},
			); err != nil {
				utils.Logger.Panicf("连接不上etcd服务器: %v", err)
			} else {
				utils.Logger.Infoln("连接etcd成功")
				serviceHub = &ServiceHub{
					client:             client,
					heartbeatFrequency: heartbeatFrequency, //租约的有效期
					loadBalancer:       &RoundRobin{},
				}
			}
		})
	}
	return serviceHub
}

func (hub *ServiceHub) Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error) {
	ctx := context.Background()
	if leaseID <= 0 {
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			utils.Logger.Warnf("创建租约失败：%v", err)
			return 0, err
		} else {
			key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
			if _, err = hub.client.Put(ctx, key, "", etcdv3.WithLease(lease.ID)); err != nil { //只需要key，不需要value
				utils.Logger.Warnf("写入服务%s对应的节点%s失败:%v", service, endpoint, err)
				return lease.ID, err
			} else {
				utils.Logger.Debugf("写入服务%s对应的节点%s成功!!!", service, endpoint)
				return lease.ID, nil
			}
		}
	} else {
		//续租
		if _, err := hub.client.KeepAliveOnce(ctx, leaseID); err == rpctypes.ErrLeaseNotFound { //续约一次，到期后还得再续约
			return hub.Regist(service, endpoint, 0) //找不到租约，走注册流程(把leaseID置为0)
		} else if err != nil {
			utils.Logger.Warnf("续约失败:%v", err)
			return 0, err
		} else {
			utils.Logger.Debugf("服务%s对应的节点%s续约成功", service, endpoint)
			return leaseID, nil
		}
	}
}

// 注销服务
func (hub *ServiceHub) UnRegist(service string, endpoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
	if _, err := hub.client.Delete(ctx, key); err != nil {
		utils.Logger.Warnf("注销服务%s对应的节点%s失败: %v", service, endpoint, err)
		return err
	} else {
		utils.Logger.Infof("注销服务%s对应的节点%s", service, endpoint)
		return nil
	}
}

func (hub *ServiceHub) GetServiceEndPoints(service string) []string {
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcdv3.WithPrefix()); err != nil {
		utils.Logger.Warnf("获取服务:%s的节点失败:%v", service, err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			path := strings.Split(string(kv.Key), "/")
			endpoints = append(endpoints, path[len(path)-1])
		}
		return endpoints
	}
}

func (hub *ServiceHub) GetServiceEndPoint(service string) string {
	return hub.loadBalancer.Take(hub.GetServiceEndPoints(service))
}
