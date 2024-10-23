package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// 起始时间戳 (2023-01-01 00:00:00 UTC)
var epoch int64 = 1672531200000

// Generator 雪花算法ID生成器
type Generator struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
	lastTime  int64 // 新增：上次生成ID的时间
}

// NewGenerator 创建一个新的 Generator 实例
func NewGenerator(node int64) (*Generator, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("节点ID超出范围")
	}
	return &Generator{
		timestamp: 0,
		node:      node,
		step:      0,
		lastTime:  0,
	}, nil
}

// Next 生成下一个唯一ID
func (g *Generator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixNano() / 1e6

	if now < g.lastTime {
		// 时钟回拨，等待直到追上lastTime
		for now <= g.lastTime {
			now = time.Now().UnixNano() / 1e6
		}
	}

	if g.timestamp == now {
		g.step++
		if g.step > stepMax {
			for now <= g.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
			g.step = 0
		}
	} else {
		g.step = 0
	}

	g.timestamp = now
	g.lastTime = now // 更新lastTime

	return (now-epoch)<<timeShift | (g.node << nodeShift) | (g.step)
}

// NextBatch 生成指定数量的唯一ID
func (g *Generator) NextBatch(count int) []int64 {
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = g.Next()
	}
	return ids
}
