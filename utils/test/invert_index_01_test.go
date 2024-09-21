package test

import (
	"fmt"
	"search_engine/utils"
	"testing"
)

func TestInvertIndex01(t *testing.T) {
	docs := []*utils.Doc{{Id: 01, Keys: []string{"go", "编程", "高并发"}}, {Id: 02, Keys: []string{"python", "编程", "大数据"}}}
	mp := utils.BuildInvertIndex(docs)
	for k, v := range mp {
		fmt.Println(k,"------", v)
	}
}
