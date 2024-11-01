package test

import (
	"fmt"
	"search_engine/index_service"
	"testing"
	"time"
)

var (
	qps         = 10
	etcdServers = []string{"127.0.0.1:2379"}
	service     = "test"

	endpoint1 = "127.0.0.1:5000"
	endpoint2 = "127.0.0.2:5000"
	endpoint3 = "127.0.0.3:5000"
)

func TestServiceHub(t *testing.T) {
	proxy := index_service.GetServiceHubProxy(etcdServers, 3, qps)

	proxy.Regist(service, endpoint1, 0)
	defer proxy.UnRegist(service, endpoint1)
	result := proxy.GetServiceEndpoints(service)
	fmt.Println(result)

	proxy.Regist(service, endpoint2, 0)
	defer proxy.UnRegist(service, endpoint2)
	result = proxy.GetServiceEndpoints(service)
	fmt.Println(result)

	proxy.Regist(service, endpoint3, 0)
	defer proxy.UnRegist(service, endpoint3)
	result = proxy.GetServiceEndpoints(service)
	fmt.Println(result)

	time.Sleep(1 * time.Second)

	for i := 0; i < qps+5; i++ {
		endpoints := proxy.GetServiceEndpoints(service)
		fmt.Println(endpoints)
	}

}
