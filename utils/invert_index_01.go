package utils

type Doc struct {
	Id   int
	Keys []string
}

// 构建初始版本的倒排索引
func BuildInvertIndex(docs []*Doc) map[string][]int {
	result := make(map[string][]int, 100)
	for _, doc := range docs {
		for _, key := range doc.Keys {
			result[key] = append(result[key], doc.Id)
		}
	}
	return result
}
