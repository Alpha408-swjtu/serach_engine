package kvdb

import (
	"errors"
	"sync/atomic"

	bolt "go.etcd.io/bbolt"
)

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

func (s *Bolt) Get(k []byte) ([]byte, error) {
	var result []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		result = tx.Bucket(s.bucket).Get(k)
		return nil
	})
	if len(result) == 0 {
		return nil, errors.New("没有数据")
	}
	return result, err
}

func (s *Bolt) BatchGet(keys [][]byte) ([][]byte, error) {
	var err error
	result := make([][]byte, 0, len(keys))
	s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			ival := tx.Bucket(s.bucket).Get(key)
			result[i] = ival
		}
		return nil
	})
	return result, err
}

func (s *Bolt) Delete(k []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Delete(k)
	})
}

func (s *Bolt) BatchDelete(keys [][]byte) error {
	var err error
	s.db.Batch(func(tx *bolt.Tx) error {
		for _, key := range keys {
			tx.Bucket(s.bucket).Delete(key)
		}
		return nil
	})
	return err
}

func (s *Bolt) Has(key []byte) bool {
	var b []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(s.bucket).Get(key)
		return nil
	})
	if err != nil || string(b) == "" {
		return false
	}
	return true
}

func (s *Bolt) IterDB(fn func(k, v []byte) error) int64 {
	var result int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			} else {
				atomic.AddInt64(&result, 1)
			}
		}
		return nil
	})
	return result
}

func (s *Bolt) IterKey(fn func(k []byte) error) int64 {
	var result int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := fn(k); err != nil {
				return err
			} else {
				atomic.AddInt64(&result, 1)
			}
		}
		return nil
	})
	return result
}

func (s *Bolt) Close() error {
	return s.db.Close()
}
