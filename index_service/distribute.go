package index_service

import (
	"context"
	"fmt"
	"search_engine/types"
	"search_engine/utils"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// 暴露给客户端的接口
type Sentinel struct {
	hub      IServiceHub
	connPool sync.Map //跟每个worker建立rpc连接,存放endpoint和对应的rpc连接
}

func NewSentinel(etcdServers []string) *Sentinel {
	return &Sentinel{
		hub:      GetServiceHubProxy(etcdServers, 3, 100),
		connPool: sync.Map{},
	}
}

func (s *Sentinel) GetGrpcConn(endpoint string) *grpc.ClientConn {
	//已有连接
	if v, exists := s.connPool.Load(endpoint); exists {
		conn := v.(*grpc.ClientConn)

		//校验链接是否可用
		if conn.GetState() == connectivity.TransientFailure || conn.GetState() == connectivity.Shutdown {
			utils.Logger.Warn("链接不可用")
			conn.Close()
			s.connPool.Delete(endpoint)
		} else {
			return conn
		}
	}

	//无，创建连接并存储
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		utils.Logger.Panicf("grpc连接不上终端%s:%v", endpoint, err)
		return nil
	}

	utils.Logger.Infof("grpc成功建立与%s的连接", endpoint)
	return conn
}

func (s *Sentinel) AddDoc(doc types.Document) (int, error) {
	endPoint := s.hub.GetServiceEndpoint(INDEX_SERVICE)
	if len(endPoint) == 0 {
		return 0, fmt.Errorf("没有合适的节点")
	}

	conn := s.GetGrpcConn(endPoint)
	if conn == nil {
		return 0, fmt.Errorf("建立到%s的连接失败", endPoint)
	}

	client := NewIndexServiceClient(conn)
	aff, err := client.AddDoc(context.Background(), &doc)
	if err != nil {
		return 0, err
	}
	return int(aff.Count), nil
}

func (s *Sentinel) DeleteDoc(docId string) int {
	//获取所有服务区，每一台都看有没有目标
	endPoints := s.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endPoints) == 0 {
		return 0
	}

	var n int32
	wg := sync.WaitGroup{}
	wg.Add(len(endPoints))

	for _, endpoint := range endPoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := s.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				affc, err := client.DeleteDoc(context.Background(), &DocId{docId})
				if err != nil {
					utils.Logger.Errorf("删除%s上的id为%s的文章有误:%v", endpoint, docId, err)

				} else {
					if affc.Count > 0 {
						atomic.AddInt32(&n, affc.Count)
						utils.Logger.Infof("删除节点%s上的id为%s的文章成功", endpoint, docId)
					}
				}
			}
		}(endpoint)

	}
	wg.Wait()
	return int(atomic.LoadInt32(&n))
}

func (s *Sentinel) Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document {
	endPoints := s.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endPoints) == 0 {
		return nil
	}

	docs := make([]*types.Document, 0, 1000)
	ch := make(chan *types.Document, 1000)

	wg := sync.WaitGroup{}
	wg.Add(len(endPoints))
	for _, endpoint := range endPoints {
		go func(string) {
			defer wg.Done()
			conn := s.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				result, err := client.Search(context.Background(), &SearchRequest{
					Query:   query,
					OnFlag:  onFlag,
					OffFlag: offFlag,
					OrFlags: orFlags,
				})
				if err != nil {
					utils.Logger.Errorf("查不到节点:%s上的文章", endpoint)
				} else {
					if len(result.Results) > 0 {
						utils.Logger.Infof("在%s上查到文章%d篇", endpoint, len(result.Results))
						for _, doc := range result.Results {
							ch <- doc
						}
					}
				}
			}
		}(endpoint)
	}

	receiveCh := make(chan struct{})
	go func() {
		for {
			result, ok := <-ch
			if !ok {
				break
			}
			docs = append(docs, result)
		}
		receiveCh <- struct{}{}
	}()
	wg.Wait()
	close(ch)
	<-receiveCh
	return docs
}
