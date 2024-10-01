package test

import (
	"fmt"
	kvdb "search_engine/internal/kv_db"
	"testing"
)

func TestBadger(t *testing.T) {
	db, err := kvdb.GetKvDb(kvdb.BADGER, "badger_db")
	if err != nil {
		fmt.Println("连接数据库失败")
		t.Fail()
	}

	defer db.Close()

	Seg()
	k := []byte("a")
	v := []byte("A")
	if err = db.Set(k, v); err != nil {
		fmt.Println("插入失败")
		t.Fail()
	}
	fmt.Println("插入成功")
	Seg()

	v, err = db.Get(k)
	if err != nil {
		fmt.Println("读取失败")
		t.Fail()
	}
	fmt.Println(string(v))

	Seg()
	if err = db.Delete(k); err != nil {
		fmt.Println("删除失败")
		t.Fail()
	}
	fmt.Println("删除成功")
}
