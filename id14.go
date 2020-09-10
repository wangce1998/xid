package xid

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const MAX = 1000

type ID14Generator struct {
	mu    sync.Mutex
	time  int64
	start int
	step  int
	node  int
}

// ---10--- | ---1--- | ---3---  // 260年，10000/s【最大时间 9999999999 => 2286-11-21 01:46:39】
// 时间戳(s) | 机器 id  |  计数器
func (g *ID14Generator) Next() int64 {
	g.mu.Lock()
	now := time.Now().Unix()
	if now == g.time {
		g.step = (g.step + 1) % MAX
		if g.step == g.start {
			for now <= g.time {
				now = time.Now().Unix()
			}
			g.start = rand.Intn(1000)
			g.step = g.start
		}
	} else {
		g.start = rand.Intn(1000)
		g.step = g.start
	}
	g.time = now
	num := fmt.Sprintf("%010d%d%03d", now, g.node, g.step)

	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		panic(err)
	}
	g.mu.Unlock()
	return n
}

func NewID14Gen(node int) (IDGen, error) {
	return &ID14Generator{node: node}, nil
}
