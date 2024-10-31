package index_service

import (
	"bytes"
	"encoding/gob"
	kvdb "search_engine/internal/kv_db"
	reverseindex "search_engine/internal/reverse_index"
	"search_engine/types"
	"search_engine/utils"
	"strings"
	"sync/atomic"
)

// 将正排索引和倒排索引放在一起
type Indexer struct {
	forwardIndex kvdb.IKeyValueDB
	reverseIndex reverseindex.IReverseIndexr
	maxIntId     uint64
}

func (indexer *Indexer) Init(DocNumEstimate int, dbType int, DataDir string) error {
	db, err := kvdb.GetKvDb(dbType, DataDir)
	if err != nil {
		return err
	}
	indexer.forwardIndex = db
	indexer.reverseIndex = reverseindex.NewSkipListReserveIndex(DocNumEstimate)
	return nil
}

func (indexer *Indexer) Close() error {
	return indexer.forwardIndex.Close()
}

func (indexer *Indexer) AddDoc(doc types.Document) (int, error) {
	docId := strings.TrimSpace(doc.Id)
	if len(docId) == 0 {
		return 0, nil
	}
	indexer.DeleteDoc(docId)

	doc.IntId = atomic.AddUint64(&indexer.maxIntId, 1)

	//写入正排索引,业务侧id和文档序列化结果
	var value bytes.Buffer
	encoder := gob.NewEncoder(&value)
	if err := encoder.Encode(doc); err == nil {
		indexer.forwardIndex.Set([]byte(docId), value.Bytes())
	} else {
		return 0, err
	}

	//写入倒排索引
	indexer.reverseIndex.Add(doc)
	return 1, nil
}

func (indexer *Indexer) DeleteDoc(docId string) int {
	n := 0
	forwardKey := []byte(docId)
	//先读正排，用key拿到文章中的keyword
	docBs, err := indexer.forwardIndex.Get(forwardKey)
	if err == nil {
		reader := bytes.NewReader([]byte{})
		if len(docBs) > 0 {
			n = 1
			reader.Reset(docBs)
			decoder := gob.NewDecoder(reader)
			var doc types.Document
			err := decoder.Decode(&doc)
			if err == nil {
				for _, keyword := range doc.Keywords {
					indexer.reverseIndex.Delete(doc.IntId, keyword)
				}
			}
		}
	}
	//从正排索引删除
	indexer.forwardIndex.Delete(forwardKey)
	return n
}

// 系统重启，直接从正排索引加载文件到倒排
func (indexer *Indexer) LoadFromIndexFile() int {
	reader := bytes.NewReader([]byte{})
	n := indexer.forwardIndex.IterDB(func(k, v []byte) error {
		reader.Reset(v)
		decoder := gob.NewDecoder(reader)
		var doc types.Document
		err := decoder.Decode(&doc)
		if err != nil {
			utils.Logger.Panic("解码失败")
			return nil
		}
		indexer.reverseIndex.Add(doc)
		return err
	})
	utils.Logger.Infof("从正排索引获取到:%d条数据", n)
	return int(n)
}

// 搜索思路：先从倒排索引搜索出文章id，再从kv数据库读取文章完整内容并解码
func (indexer *Indexer) Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document {
	docIds := indexer.reverseIndex.Search(query, onFlag, offFlag, orFlags)
	if len(docIds) == 0 {
		return nil
	}
	keys := make([][]byte, 0, len(docIds))
	for _, docId := range docIds {
		keys = append(keys, []byte(docId))
	}
	data, err := indexer.forwardIndex.BatchGet(keys)
	if err != nil {
		utils.Logger.Warn("读取kv数据库失败")
		return nil
	}
	result := make([]*types.Document, 0, len(data))
	reader := bytes.NewReader([]byte{})
	for _, docBs := range data {
		if len(docBs) > 0 {
			reader.Reset(docBs)
			decoder := gob.NewDecoder(reader)
			var doc types.Document
			err := decoder.Decode(&doc)
			if err != nil {
				result = append(result, &doc)
			}
		}
	}
	return result
}
