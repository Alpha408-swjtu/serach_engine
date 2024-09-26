package reverseindex

import (
	"runtime"
	"search_engine/types"
	"search_engine/utils"
	"sync"

	"github.com/huandu/skiplist"
	farmhash "github.com/leemcloughlin/gofarmhash"
)

// 基于跳表实现倒排索引
type SkipListReserveIndex struct {
	table *utils.ConcurrentHashMap
	locks []sync.RWMutex
}

// 定义跳表map中value的类型
type SkipListValue struct {
	Id          string
	BitsFeature uint64
}

func NewSkipListReserveIndex(DocsNum int) *SkipListReserveIndex {
	result := new(SkipListReserveIndex)
	result.table = utils.NewConcurrentMap(runtime.NumCPU(), DocsNum) //假设key的值和文章数量相当
	result.locks = make([]sync.RWMutex, 1000)
	return result
}

func (indexer *SkipListReserveIndex) getLock(key string) *sync.RWMutex {
	h := int(farmhash.Hash32WithSeed([]byte(key), 0))
	return &indexer.locks[h%len(indexer.locks)]
}

func (indexer *SkipListReserveIndex) Add(doc *types.Document) {
	//遍历文章中的所有关键词
	for _, Keyword := range doc.Keywords {
		key := Keyword.Tostring()
		lock := indexer.getLock(key)
		lock.Lock()
		defer lock.Unlock()
		skipListValue := &SkipListValue{doc.Id, doc.BitsFeature}
		if value, exists := indexer.table.Get(key); exists {
			list := value.(*skiplist.SkipList)
			list.Set(doc.Id, skipListValue)
		} else {
			list := skiplist.New(skiplist.Uint64)
			list.Set(doc.Id, skipListValue)
		}
	}
}

func (indexer *SkipListReserveIndex) Delete(IntId uint64, keyWord *types.Keyword) {
	key := keyWord.Tostring()
	lock := indexer.getLock(key)
	lock.Lock()
	defer lock.Unlock()
	if value, exists := indexer.table.Get(key); exists {
		list := value.(*skiplist.SkipList)
		list.Remove(IntId)
	}
}

// 多个跳表求交集
func IntersectionOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}
	result := skiplist.New(skiplist.Uint64)
	currNodes := make([]*skiplist.Element, len(lists)) //每一个表都分配一个指针
	for i, list := range lists {
		if list == nil || list.Len() == 0 {
			return nil //有一条空表，则交集为空
		}
		currNodes[i] = list.Front() //初始化时指针都指向调表的首元素
	}

	for {
		maxList := make(map[int]struct{}, len(currNodes))
		var maxValue uint64
		for i, node := range currNodes {
			if node.Key().(uint64) > maxValue {
				maxValue = node.Key().(uint64)
				maxList = map[int]struct{}{i: {}}
			} else if node.Key().(uint64) == maxValue {
				maxList[i] = struct{}{}
			}
		}
		if len(maxList) == len(currNodes) {
			result.Set(currNodes[0].Key(), currNodes[0].Value)
			for i, node := range currNodes {
				currNodes[i] = node.Next()
				if currNodes[i] == nil {
					return result
				}
			}
		} else {
			for i, node := range currNodes {
				if _, exists := maxList[i]; !exists {
					currNodes[i] = node.Next()
					if currNodes[i] == nil {
						return result
					}
				}
			}
		}
	}
}

// 多个跳表求并集
func UnionsetOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}
	result := skiplist.New(skiplist.Uint64)
	keySet := make(map[any]struct{}, 1000)
	for _, list := range lists {
		if list == nil {
			continue
		}
		node := list.Front()
		for node != nil {
			if _, exists := keySet[node.Key()]; !exists {
				result.Set(node.Key(), node.Value)
				keySet[node.Key()] = struct{}{}
			}
			node = node.Next()
		}
	}
	return result
}
