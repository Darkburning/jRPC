package service

import (
	"math"
)

// 定义服务，服务必须包含入参和出参

func Add(a, b int) int {
	return a + b
}

func Substract(a, b int) int {
	return a - b
}

func Consub(a, b int) int {
	return b - a
}

func Multi(a, b int) int {
	return a * b
}

func Divide(a, b int) float64 {
	return float64(a / b)
}

func Condiv(a, b int) float64 {
	return float64(b / a)
}

func Power(a, b int) int {
	res := 1
	for i := 0; i < b; i++ {
		res *= a
	}
	return res
}

func Mod(a, b int) int {
	return a % b
}

func Sqrtmul(a, b int) float64 {
	return math.Sqrt(float64(a)) * math.Sqrt(float64(b))
}

func Triangle(a, b int) float64 {
	return float64(a * b / 2)
}
