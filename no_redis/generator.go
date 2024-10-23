package no_redis

import (
	"sync/atomic"
	"time"
)

// Generator 是生成唯一数字的主要结构体
type Generator struct {
	counter int64
	prefix  int64
}

// NewGenerator 创建一个新的 Generator 实例
func NewGenerator() *Generator {
	return &Generator{
		counter: 0,
		prefix:  time.Now().UnixNano() / 1e6, // 使用毫秒级时间戳作为前缀
	}
}

// Next 生成下一个唯一数字
func (g *Generator) Next() int64 {
	count := atomic.AddInt64(&g.counter, 1)
	return (g.prefix << 22) | (count & 0x3FFFFF) // 组合前缀和计数器
}

// NextBatch 生成指定数量的唯一数字
func (g *Generator) NextBatch(count int) []int64 {
	numbers := make([]int64, count)
	for i := 0; i < count; i++ {
		numbers[i] = g.Next()
	}
	return numbers
}
