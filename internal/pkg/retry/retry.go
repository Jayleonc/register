package retry

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Async 重试任务或异步任务（非定时任务，定时任务请用 job）
type Async interface {
	AsyncFunc(ctx context.Context) error
	NextInterval() time.Duration
}

// Scheduler 负责调度和执行任务
type Scheduler struct {
	tasks   []Async
	ctx     context.Context
	cancel  func()
	mu      sync.Mutex
	started bool
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		ctx:    ctx,
		cancel: cancel,
	}
	go func() {
		s.Start()
	}()
	return s
}

// RegisterAsync 用于注册任务到调度器
func (s *Scheduler) RegisterAsync(task ...Async) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		// 如果调度器已经启动，立即开始执行新注册的任务
		// 防止调用方在运行时的某个时机注册任务
		for _, t := range task {
			go s.execute(t)
		}
	}
	s.tasks = append(s.tasks, task...)
}

// Start 启动调度器，遍历并执行每个任务
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return
	}
	s.started = true

	for _, task := range s.tasks {
		go s.execute(task)
	}
}

func (s *Scheduler) execute(t Async) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if err := t.AsyncFunc(s.ctx); err != nil {
				fmt.Println("任务执行出错:", err)
			}

			nextInterval := t.NextInterval()
			if nextInterval > 0 {
				select {
				case <-s.ctx.Done():
					return
				case <-time.After(nextInterval):
				}
			}
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cancel()
}
