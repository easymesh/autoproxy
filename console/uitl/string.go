package util

import (
	"fmt"
	"strconv"
	"strings"
)

func Atoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

func ListAtoi(ss []string) []int {
	var output []int
	for _, v := range ss {
		output = append(output, Atoi(v))
	}
	return output
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func StringCat(v, pre, suf string) string {
	begin := strings.Index(v, pre)
	if begin == -1 {
		return ""
	}
	begin += len(pre)
	end := strings.Index(v[begin:], suf)
	if end == -1 {
		return ""
	}
	return v[begin: begin + end]
}