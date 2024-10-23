package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/nextnum"
)

func main() {
	// 创建一个新的生成器实例
	generator := nextnum.NewGenerator("localhost:6379", "order_counter")
	defer generator.Close()

	ctx := context.Background()

	// 生成单个数字
	num, err := generator.Next(ctx)
	if err != nil {
		log.Fatalf("生成单个数字时出错: %v", err)
	}
	fmt.Printf("生成的单个数字: %d\n", num)

	// 生成多个数字
	count := 5
	numbers, err := generator.NextBatch(ctx, count)
	if err != nil {
		log.Fatalf("生成多个数字时出错: %v", err)
	}
	fmt.Printf("生成的 %d 个数字: %v\n", count, numbers)
}
