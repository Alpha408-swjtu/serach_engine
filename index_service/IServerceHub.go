package index_service

import etcdv3 "go.etcd.io/etcd/client/v3"

type IServiceHub interface {
	Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error)
	UnRegist(service string, endpoint string) error
	GetServiceEndpoints(service string) []string
	GetServiceEndpoint(service string) string
}
