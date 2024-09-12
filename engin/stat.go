package engin

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	total_session int64
	total_size    int64
	version       string
)

const (
	ONE_GB float64 = 1024 * 1024 * 1024
	ONE_MB float64 = 1024 * 1024
	ONE_KB float64 = 1024
)

func StatUpdate(session int64, size int64) {
	atomic.AddInt64(&total_session, session)
	atomic.AddInt64(&total_size, size)
}

func display_size() string {
	size := float64(total_size)
	if size > ONE_GB {
		return fmt.Sprintf("%.2f GB", size/ONE_GB)
	}
	if size > ONE_MB {
		return fmt.Sprintf("%.2f MB", size/ONE_MB)
	}
	if size > ONE_KB {
		return fmt.Sprintf("%.2f KB", size/ONE_KB)
	}
	return fmt.Sprintf("%.f B", size)
}

func Display() {
	for {
		time.Sleep(time.Second)
		fmt.Print("\r                                                             ")
		fmt.Printf("\r [version: %s session: %d, throughput: %s]", version, total_session, display_size())
	}
}

func SetVersion(v string) {
	version = v
}
