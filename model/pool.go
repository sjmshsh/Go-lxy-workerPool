package model

import (
	"fmt"
	"github.com/sjmshsh/pkg/constance"
	"github.com/sjmshsh/pkg/er"
	"sync"
)

type Pool struct {
	PreAlloc bool // 是否在创建pool的时候就预创建workers, 默认值为false
	// 当pool满的时候，新的Schedule调用是否阻塞当前goroutine。默认值：true
	// 如果block = false，则Schedule返回ErrNoWorkerAvailInPool
	Block    bool
	capacity int            // workerpool大小
	active   chan struct{}  // 对应架构图中的active channel
	tasks    chan Task      // 对应架构图中的task channel
	wg       sync.WaitGroup // 用于在pool销毁时等待所有worker退出
	quit     chan struct{}  // 用于通知各个worker退出的信号channel
}

// New 用于创建一个pool类型实例，并将pool池的worker管理机制运行起来
func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = constance.DefaultCapacity
	}
	if capacity > constance.MaxCapacity {
		capacity = constance.MaxCapacity
	}
	p := &Pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
	}
	go p.run()
	return p
}

// run run方法内是一个无限循环, 循环体中使用select监视Pool类型实例的两个channel
// quit 和 active。当接收到来自quit channel的退出信号的时候，这个Goroutine会结束运行
// 而当active channel可写的时候，run方法就会创建一个新的worker Goroutine
func (p *Pool) run() {
	idx := 0
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			// create a new worker
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: starter\n", i)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return er.ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

func (p *Pool) Free() {

}
