package kvdb

import (
	"fmt"
	"os"
	"strings"
)

//用户可以自行选择数据库

const (
	BOLT = iota
	BADGER
)

type IKeyValueDB interface {
	Open() error
	GetDbPath() string
	Set(k, v []byte) error
	BatchSet(keys, values [][]byte) error //批量写入k和v
	Get(k []byte) ([]byte, error)
	BatchGet(keys [][]byte) ([][]byte, error)
	Delete(k []byte) error
	BatchDelete(keys [][]byte) error
	Has(key []byte) bool                     //判断是否有key
	IterDB(fn func(k, v []byte) error) int64 //遍历数据库，返回数据条数
	IterKey(fn func(k []byte) error) int64   //遍历所有key，返回数据条数
	Close() error                            //关闭数据库
}

func GetKvDb(dbtype int, path string) (IKeyValueDB, error) {
	paths := strings.Split(path, "/")
	parentPath := strings.Join(paths[0:len(paths)-1], "/") //排除最后一个路径，只取得父目录

	info, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		fmt.Printf("创建路径:%s", parentPath)
		os.MkdirAll(parentPath, 0o600)
	} else {
		if info.Mode().IsRegular() {
			fmt.Printf("%s 是一个文件，不是目录，需要删除。", parentPath)
			os.Remove(parentPath)
		}
	}

	var db IKeyValueDB
	switch dbtype {
	case BADGER:
		db = new(Badger).WithDataPath(path)
	default:
		db = new(Bolt).WithDataPath(path).WithBucket("radic")
	}

	err = db.Open()

	return db, err
}
