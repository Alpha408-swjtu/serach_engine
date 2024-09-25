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
