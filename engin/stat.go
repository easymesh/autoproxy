package engin

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	total_session int64
	total_size    int64
)

func StatUpdate(session int64, size int64) {
	atomic.AddInt64(&total_session, session)
	atomic.AddInt64(&total_size, size)
}

func display() {
	for {
		time.Sleep(time.Second)
		fmt.Printf("\r                                                                           ")
		fmt.Printf("\r session: %d, speed: 10kb/s", total_session)
	}
}

func init() {
	go display()
}
