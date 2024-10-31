package test

import (
	"fmt"
	"search_engine/index_service"
	"testing"
)

var (
	endpoints = []string{"127.0.0.1:2379"}
	service   = "test"
)

func TestServiceHub(t *testing.T) {
	hub := index_service.GetServiceHub(endpoints, 1)
	hub.Regist(service, "127.0.0.1:2379", 0)
	s := hub.GetServiceEndPoints(service)
	fmt.Println(s)
	defer hub.UnRegist(service, "127.0.0.1:2379")

}
