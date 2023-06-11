package main

import "time"

// 定义服务，服务必须包含入参和出参

func Sum(args ...float64) float64 {
	sum := 0.0
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return sum
}

func Product(args ...float64) float64 {
	product := 1.0
	for i := 0; i < len(args); i++ {
		product *= args[i]
	}
	return product
}

func Sleep(seconds float64) float64 {
	time.Sleep(time.Duration(seconds) * time.Second)
	return seconds
}
