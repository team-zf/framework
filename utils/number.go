package utils

import (
	"math/rand"
	"time"
)

// 指定范围内随机生成一个值 (含max)
func Range(args ...int) int {
	rand.Seed(time.Now().UnixNano())
	switch len(args) {
	// 0 ~ max
	case 1:
		return rand.Intn(args[0] + 1)
	// min ~ max
	case 2:
		return rand.Intn(args[1]+1-args[0]) + args[0]
	}
	return 0
}

// 万分比概率
func Percent(v int) bool {
	return Range(1, 10000) <= v
}

// 万分比概率值
func PercentV() int {
	return Range(1, 10000)
}

// 取出最小的值
func Min(args ...int) int {
	var val int
	for _, n := range args {
		if n < val {
			val = n
		}
	}
	return val
}

// 取出最小的值
func Min64(args ...int64) int64 {
	var val int64
	for _, n := range args {
		if n < val {
			val = n
		}
	}
	return val
}

// 取出最大的值
func Max(args ...int) int {
	val := 0
	for _, n := range args {
		if n > val {
			val = n
		}
	}
	return val
}
