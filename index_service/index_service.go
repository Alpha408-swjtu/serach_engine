package index_service

import (
	"context"
	types "search_engine/types"
	"search_engine/utils"
	"strconv"
	"time"
)

const (
	INDEX_SERVICE = "index_service"
)

type IndexServiceWorker struct {
	Indexer  *Indexer //正排倒排放一块
	hub      *ServiceHub
	selfAddr string
}

func (service *IndexServiceWorker) Init(DocNumEstimate int, dbtype int, DataDir string, etcdServers []string, servicePort int) error {
	//开启正排索引
	service.Indexer = new(Indexer)
	service.Indexer.Init(DocNumEstimate, dbtype, DataDir)

	//注册
	if len(etcdServers) > 0 {
		if servicePort <= 1024 {
			utils.Logger.Errorf("端口非法:%v", servicePort)
		}
		selfLocalIP, err := utils.GetLocalIP()
		if err != nil {
			utils.Logger.Panicf("获取不到IP地址:%v", err)
			return err
		}

		selfLocalIP = "127.0.0.1" //单机演示写死，多机分布再改
		service.selfAddr = selfLocalIP + ":" + strconv.Itoa(servicePort)

		var heartBeat int64 = 3
		hub := GetServiceHub(etcdServers, heartBeat)
		leaseId, err := hub.Regist(INDEX_SERVICE, service.selfAddr, 0)
		if err != nil {
			utils.Logger.Panicf("创建租约失败:%v", err)
			return err
		}
		service.hub = hub

		go func() {
			for { //一直注册
				hub.Regist(INDEX_SERVICE, service.selfAddr, leaseId)
				time.Sleep(time.Duration(heartBeat)*time.Second - 100*time.Millisecond)
			}
		}()
	}

	return nil
}

func (service *IndexServiceWorker) LoadFromIndexFile() int {
	return service.Indexer.LoadFromIndexFile()
}

func (service *IndexServiceWorker) Close() error {
	//注销etcd
	if service.hub != nil {
		service.hub.UnRegist(INDEX_SERVICE, service.selfAddr)
	}
	//关闭正排索引(kvdb)
	return service.Indexer.Close()
}

func (service *IndexServiceWorker) DeleteDoc(ctx context.Context, docId *DocId) (*AffectedCount, error) {
	return &AffectedCount{int32(service.Indexer.DeleteDoc(docId.DocId))}, nil
}

// 向索引中添加文档(如果已存在，会先删除)
func (service *IndexServiceWorker) AddDoc(ctx context.Context, doc *types.Document) (*AffectedCount, error) {
	n, err := service.Indexer.AddDoc(*doc)
	return &AffectedCount{int32(n)}, err
}

// 检索，返回文档列表
func (service *IndexServiceWorker) Search(ctx context.Context, request *SearchRequest) (*SearchResult, error) {
	result := service.Indexer.Search(request.Query, request.OnFlag, request.OffFlag, request.OrFlags)
	return &SearchResult{Results: result}, nil
}
