package reverseindex

import "search_engine/types"

//封装倒排索引必须实现的功能
type IReverseIndexr interface {
	Add(doc types.Document)                                                                  //添加文章
	Delete(IntId uint64, keyWord *types.Keyword)                                             //删除文章
	Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []string //搜索文章
}
