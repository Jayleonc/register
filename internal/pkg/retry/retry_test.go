package retry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockTask struct {
	ExecuteCalled       bool
	ExecutionTimestamps []time.Time
	NextIntervalVal     time.Duration
}

func (m *MockTask) AsyncFunc(ctx context.Context) error {
	m.ExecutionTimestamps = append(m.ExecutionTimestamps, time.Now())
	m.ExecuteCalled = true
	return nil
}

func (m *MockTask) NextInterval() time.Duration {
	return m.NextIntervalVal
}

func TestAsyncScheduler(t *testing.T) {
	task := &MockTask{NextIntervalVal: 500 * time.Millisecond}
	scheduler := NewScheduler()
	scheduler.RegisterAsync(task)

	time.Sleep(2 * time.Second)

	assert.True(t, task.ExecuteCalled)

	// 测试停止调度器
	scheduler.Stop()

	// 验证调度器停止后任务不再执行
	task.ExecuteCalled = false
	time.Sleep(2 * time.Second)

	assert.True(t, !task.ExecuteCalled)
}

func TestTaskExecutionInterval(t *testing.T) {
	task := &MockTask{NextIntervalVal: 10 * time.Millisecond}
	scheduler := NewScheduler()
	scheduler.RegisterAsync(task)

	time.Sleep(100 * time.Millisecond)

	scheduler.Stop()

	// 验证任务至少执行了两次
	assert.True(t, len(task.ExecutionTimestamps) > 5)

	// 验证执行间隔
	executionInterval := task.ExecutionTimestamps[1].Sub(task.ExecutionTimestamps[0])
	assert.True(t, executionInterval >= task.NextIntervalVal)
}
