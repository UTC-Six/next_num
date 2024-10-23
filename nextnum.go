package nextnum

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

// Generator 是生成唯一数字的主要结构体
type Generator struct {
	redisClient *redis.Client
	mutex       sync.Mutex
	counterKey  string
}

// NewGenerator 创建一个新的 Generator 实例
func NewGenerator(redisAddr, counterKey string) *Generator {
	return &Generator{
		redisClient: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
		counterKey: counterKey,
	}
}

// Next 生成下一个唯一数字
func (g *Generator) Next(ctx context.Context) (int64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	num, err := g.redisClient.Incr(ctx, g.counterKey).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

// NextBatch 生成指定数量的唯一数字
func (g *Generator) NextBatch(ctx context.Context, count int) ([]int64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	pipe := g.redisClient.Pipeline()
	for i := 0; i < count; i++ {
		pipe.Incr(ctx, g.counterKey)
	}
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	numbers := make([]int64, count)
	for i, cmder := range cmders {
		num, err := cmder.(*redis.IntCmd).Result()
		if err != nil {
			return nil, err
		}
		numbers[i] = num
	}
	return numbers, nil
}

// Close 关闭 Redis 连接
func (g *Generator) Close() error {
	return g.redisClient.Close()
}
