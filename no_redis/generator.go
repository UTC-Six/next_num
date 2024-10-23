package no_redis

import (
	"sync"
	"sync/atomic"
)

// Generator 是生成唯一数字的主要结构体
type Generator struct {
	counter int64
	mutex   sync.Mutex
}

// NewGenerator 创建一个新的 Generator 实例
func NewGenerator(initialValue int64) *Generator {
	return &Generator{
		counter: initialValue,
	}
}

// Next 生成下一个唯一数字
func (g *Generator) Next() int64 {
	return atomic.AddInt64(&g.counter, 1)
}

// NextBatch 生成指定数量的唯一数字
func (g *Generator) NextBatch(count int) []int64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	numbers := make([]int64, count)
	for i := 0; i < count; i++ {
		numbers[i] = atomic.AddInt64(&g.counter, 1)
	}
	return numbers
}
