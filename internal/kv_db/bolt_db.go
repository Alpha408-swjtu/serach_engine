package kvdb

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var NoDataErr = errors.New("no data")

type Bolt struct {
	db     *bolt.DB
	path   string //存储路径
	bucket []byte //表名
}

func (s *Bolt) WithDataPath(path string) *Bolt {
	s.path = path
	return s
}

func (s *Bolt) WithBucket(bucket string) *Bolt {
	s.bucket = []byte(bucket)
	return s
}

func (s *Bolt) Open() error {
	DataDir := s.GetDbPath()
	db, err := bolt.Open(DataDir, 0o600, bolt.DefaultOptions)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.bucket)
		return err
	})
	if err != nil {
		db.Close()
		return err
	} else {
		s.db = db
		return nil
	}
}

func (s *Bolt) GetDbPath() string {
	return s.path
}

func (s *Bolt) Set(k, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Put(k, v)
	})
}

func (s *Bolt) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("k和v的长度不匹配")
	}
	var err error
	s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			value := values[i]
			tx.Bucket(s.bucket).Put(key, value)
		}
		return nil
	})
	return err

}
