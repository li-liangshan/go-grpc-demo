package util

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10
	numberBits uint8 = 12
	workerMax  int64 = -1 ^ (-1 << workerBits)
	numberMax  int64 = -1 ^ (-1 << numberBits)

	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1525705533000 // 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (worker *Worker) GetId() int64 {
	worker.mu.Lock()
	defer worker.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if worker.timestamp == now {
		worker.number++
		if worker.number > numberMax {
			for now <= worker.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		worker.number = 0
		worker.timestamp = now
	}

	ID := int64((now-startTime)<<timeShift | (worker.workerId << workerShift) | (worker.number))
	return ID
}

func RunId() {
	// 生成节点实例
	node, err := NewWorker(1)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println(node.GetId())
	}
}
