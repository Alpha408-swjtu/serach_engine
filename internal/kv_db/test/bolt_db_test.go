package test

import (
	"fmt"
	kvdb "search_engine/internal/kv_db"
	"testing"
	//bolt "go.etcd.io/bbolt"
)

func Seg() {
	fmt.Println("--------------------------")
}

func TestBolt(t *testing.T) {
	db, err := kvdb.GetKvDb(kvdb.BOLT, "bolt_db")
	//defer db.Close()
	if err != nil {
		fmt.Println("连不上数据库")
		return
	}

	k := []byte("key")
	v := []byte("value")
	err = db.Set(k, v)
	if err != nil {
		fmt.Println("写入失败")
		return
	}
	Seg()
	value, err := db.Get(k)
	if err != nil {
		fmt.Println("读取失败")
		return
	}
	fmt.Println(string(value))
	Seg()
	if err = db.Delete(k); err != nil {
		fmt.Println("删除失败")
	}

	Seg()
	if err = db.Close(); err == nil {
		fmt.Println("关闭数据库")
	}
}
