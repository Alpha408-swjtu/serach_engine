package utils

import (
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

var TotalQuery int32

func Handler() {
	atomic.AddInt32(&TotalQuery, 1)
	time.Sleep(50 * time.Millisecond)
}

func CallHandler() {
	//创建一个限流器
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 10)
	n := 3
	for {
		// if limiter.AllowN(time.Now(), n) {
		// 	Handler()
		// }
		reserve := limiter.ReserveN(time.Now(), n)
		time.Sleep(reserve.Delay())
		Handler()
	}
}
