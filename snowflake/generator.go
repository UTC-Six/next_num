package snowflake

import (
	"errors"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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
	lastTime  int64
	filename  string // 用于存储lastTime的文件名
}

// NewGenerator 创建一个新的 Generator 实例
func NewGenerator(node int64, filename string) (*Generator, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("节点ID超出范围")
	}

	g := &Generator{
		node:     node,
		filename: filename,
	}

	if err := g.loadLastTime(); err != nil {
		return nil, err
	}

	// 设置退出时的保存操作
	g.setupCleanup()

	return g, nil
}

// loadLastTime 从文件加载上次的时间戳
func (g *Generator) loadLastTime() error {
	data, err := ioutil.ReadFile(g.filename)
	if err != nil {
		if os.IsNotExist(err) {
			g.lastTime = time.Now().UnixNano() / 1e6
			return nil
		}
		return err
	}

	g.lastTime, err = strconv.ParseInt(string(data), 10, 64)
	return err
}

// saveLastTime 保存最后的时间戳到文件
func (g *Generator) saveLastTime() error {
	return ioutil.WriteFile(g.filename, []byte(strconv.FormatInt(g.lastTime, 10)), 0644)
}

// setupCleanup 设置清理函数，在程序退出时保存lastTime
func (g *Generator) setupCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		g.saveLastTime()
		os.Exit(0)
	}()
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
	g.lastTime = now

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

// SaveLastTime 手动保存最后的时间戳到文件
func (g *Generator) SaveLastTime() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.saveLastTime()
}
