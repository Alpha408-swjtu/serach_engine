package utils

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	single *gorm.DB
	lock   = &sync.Mutex{}
	once   = &sync.Once{}
)

func GetDB1() *gorm.DB {
	if single == nil {
		lock.Lock()
		if single == nil {
			single, _ = gorm.Open(mysql.Open(""))
		} else {
			Logger.Info("单例已经被创建")
		}
	} else {
		Logger.Info("单例已经被创建")
	}
	return single
}

func GetDB2() *gorm.DB {
	if single == nil {
		once.Do(func() {
			single, _ = gorm.Open(mysql.Open(""))
		})
	} else {
		Logger.Info("单例已经被创建")
	}
	return single
}
