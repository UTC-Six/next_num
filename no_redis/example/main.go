package main

import (
	"fmt"
	"time"

	"github.com/yourusername/nextnum/no_redis"
)

func main() {
	// 使用当前时间戳作为初始值
	initialValue := time.Now().UnixNano()
	generator := no_redis.NewGenerator(initialValue)

	// 生成单个数字
	num := generator.Next()
	fmt.Printf("生成的单个数字: %d\n", num)

	// 生成多个数字
	count := 5
	numbers := generator.NextBatch(count)
	fmt.Printf("生成的 %d 个数字: %v\n", count, numbers)

	// 演示在不同时间生成的数字
	time.Sleep(time.Second)
	anotherGenerator := no_redis.NewGenerator(time.Now().UnixNano())
	anotherNum := anotherGenerator.Next()
	fmt.Printf("1秒后生成的数字: %d\n", anotherNum)
}