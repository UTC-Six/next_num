package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/nextnum/snowflake"
)

func main() {
	// 创建一个新的生成器实例,节点ID为1，文件名为"snowflake_lasttime.txt"
	generator, err := snowflake.NewGenerator(1, "../../snowflake_lasttime.txt")
	if err != nil {
		log.Fatalf("创建生成器失败: %v", err)
	}

	// 设置 panic 恢复和 lastTime 保存
	defer func() {
		if r := recover(); r != nil {
			generator.SaveLastTime() // 注意：我们需要在 Generator 结构体中添加这个方法
			panic(r)                 // 重新抛出 panic
		}
	}()

	// 生成单个ID
	id := generator.Next()
	fmt.Printf("生成的单个ID: %d\n", id)

	// 生成多个ID
	count := 5
	ids := generator.NextBatch(count)
	fmt.Printf("生成的 %d 个ID: %v\n", count, ids)

	// 演示在不同时间生成的ID
	time.Sleep(time.Second)
	anotherId := generator.Next()
	fmt.Printf("1秒后生成的ID: %d\n", anotherId)

	// 解析ID
	timestamp := (id >> 22) + 1672531200000 // 加回epoch
	nodeId := (id >> 12) & 0x3FF
	step := id & 0xFFF
	fmt.Printf("解析ID %d:\n", id)
	fmt.Printf("  时间戳: %s\n", time.Unix(timestamp/1000, (timestamp%1000)*1e6).Format(time.RFC3339Nano))
	fmt.Printf("  节点ID: %d\n", nodeId)
	fmt.Printf("  步骤: %d\n", step)
	panic("人造危机")
}
