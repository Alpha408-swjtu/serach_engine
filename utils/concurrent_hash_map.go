package utils

import (
	"sync"

	maps "golang.org/x/exp/maps"

	farmhash "github.com/leemcloughlin/gofarmhash"
)

// 自建并发安全的map
type ConcurrentHashMap struct {
	Mps   []map[string]any
	Locks []sync.RWMutex
	Seg   int
	Seed  uint32
}

func NewConcurrentMap(seg, cap int) *ConcurrentHashMap {
	mps := make([]map[string]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := 0; i < seg; i++ {
		mps[i] = make(map[string]any, cap/seg)
	}
	return &ConcurrentHashMap{
		Mps:   mps,
		Locks: locks,
		Seg:   seg,
		Seed:  0,
	}
}

// 判断key该写入哪个小map
func (m *ConcurrentHashMap) getSegIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), m.Seed))
	return hash % m.Seg
}

// 写入key和value
func (m *ConcurrentHashMap) Set(key string, value any) {
	index := m.getSegIndex(key)
	m.Locks[index].Lock()
	defer m.Locks[index].Unlock()
	m.Mps[index][key] = value
}

// 根据key读取value
func (m *ConcurrentHashMap) Get(key string) (any, bool) {
	index := m.getSegIndex(key)
	m.Locks[index].RLock()
	defer m.Locks[index].RUnlock()
	value, exist := m.Mps[index][key]
	return value, exist
}

// 迭代器模式：构建Next()方法可以遍历map中的元素
type MapIterator interface {
	Next() *MapEntry
}

type ConcurrentHashMapIterator struct {
	Cmp      *ConcurrentHashMap
	Keys     [][]string
	RowIndex int
	ColIndex int
}

func (m *ConcurrentHashMap) CreateIterator() *ConcurrentHashMapIterator {
	keys := make([][]string, 0, len(m.Mps))
	for _, mp := range m.Mps {
		row := maps.Keys(mp) //关键步骤，三方库获取mp里面的key
		keys = append(keys, row)
	}
	return &ConcurrentHashMapIterator{
		Cmp:      m,
		Keys:     keys,
		RowIndex: 0,
		ColIndex: 0,
	}
}

type MapEntry struct {
	Key   string
	Value any
}

func (iter *ConcurrentHashMapIterator) Next() *MapEntry {
	if iter.RowIndex >= len(iter.Keys) {
		return nil
	}

	row := iter.Keys[iter.RowIndex]
	if len(row) == 0 {
		iter.RowIndex++
		return iter.Next()
	}
	key := row[iter.ColIndex]
	value, _ := iter.Cmp.Get(key)
	if iter.ColIndex >= len(row)-1 {
		iter.ColIndex = 0
		iter.RowIndex++
	} else {
		iter.ColIndex++
	}
	return &MapEntry{Key: key, Value: value}
}
