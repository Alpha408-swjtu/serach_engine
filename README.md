# 分布式搜索引擎设计

## 文件结构

```shell
| go.mod
| go.sum
|
|--utils  
|      |--invert_index_01.go
|      |--concurrent_hash_map.go
|      |--bits.go
|      |--skip_list.go
|
```

### utils

存储项目所需要的公共用具

#### 1.invert.index_01.go

初步构建倒排索引算法，根据文档关键字获取文档ID

```go
func BuildInvertIndex(docs []*Doc) map[string][]int {
	result := make(map[string][]int, 100)
	for _, doc := range docs {
		for _, key := range doc.Keys {
			result[key] = append(result[key], doc.Id)
		}
	}
	return result
}
```

#### 2.concurrent_hash_map.go

自建并发安全的hashmap，可读写，并绑定next方法用于遍历

```go
type ConcurrentHashMap struct {
	Mps   []map[string]any
	Locks []sync.RWMutex
	Seg   int
	Seed  uint32
}


...
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

```

#### 3.bits.go

位运算相关算法

```go
...
func IsBit1(n uint64, i int) bool {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	if n&c == c {
		return true
	} else {
		return false
	}
}

func SetBit1(n uint64, i int) uint64 {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	return n | c
}

func CountBit1(n uint64) int {
	c := uint64(1)
	sum := 0
	for i := 0; i < 64; i++ {
		if c&n == c {
			sum++
		}
		c = c << 1
	}
	return sum
}
...
```

#### 4.skip_list.go

跳表相关。多跳表求交集并集。

```go
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
...//求交集的省略
```

