package service

import (
	"strings"
	"time"
)

// 定义服务，服务必须包含入参和出参

func Sum(arg1, arg2 int) int {
	return arg1 + arg2
}

func Subtract(arg1, arg2 float64) float64 {
	return arg1 - arg2
}

func Product(arg1, arg2 float64) float64 {
	return arg1 * arg2
}

func Division(arg1, arg2 float64) float64 {
	if arg2 == 0 {
		return 0
	}
	return arg1 / arg2
}

func Square(arg float64) float64 {
	return arg * arg
}

func Cube(arg float64) float64 {
	return arg * arg * arg
}

func Sleep(seconds float64) float64 {
	time.Sleep(time.Duration(seconds) * time.Second)
	return seconds
}

func Upper(str string) string {
	return strings.ToUpper(str)
}

func Lower(str string) string {
	return strings.ToLower(str)
}

func Revert(str string) string {
	r := []rune(str)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
