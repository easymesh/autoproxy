package util

import "math"

func TenThousandPercent(total int, used int) int {
	if used < 0 {
		return 0
	}
	if used > total {
		return 10000
	}
	if ( used > math.MaxInt64 / 10000 ) {
		return used / (total/10000)
	}
	return (used * 10000) / total
}

func Percent(total int, used int) int {
	if used < 0 {
		return 0
	}
	if used > total {
		return 100
	}
	if ( used > math.MaxInt64 / 100 ) {
		return used / (total/100)
	}
	return (used * 100) / total
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}