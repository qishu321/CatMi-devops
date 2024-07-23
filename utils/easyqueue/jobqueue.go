package easyqueue

import (
	"sync/atomic"
	"time"
)

type Queue interface {
	Push(JobInterface)
	Pop() JobInterface
	PopTimeout(duration time.Duration) JobInterface
}

// 单队列
type jobQueue struct {
	ch chan JobInterface
}

func createJobQueue(cap int) *jobQueue {
	return &jobQueue{
		ch: make(chan JobInterface, cap),
	}
}

// 将任务放入队列
func (jq *jobQueue) Push(jb JobInterface) {

	select {
	case jq.ch <- jb:
	default: // 说明满了
		j := jb.(*job)
		j.err = ErrOverFlow
		close(j.ch)
	}

}

func (jq *jobQueue) Pop() JobInterface {
	return <-jq.ch
}

func (jq *jobQueue) PopTimeout(duration time.Duration) JobInterface {
	select {
	case job := <-jq.ch:
		return job
	case <-time.After(duration):
		return nil
	}
}

// 多队列
type multiJobQueue struct {
	queues   []*jobQueue
	parition int

	pushIdx *idGenerator
	popIdx  *idGenerator
}

// partition * perCap = All Capacity
func createMultiJobQueue(partition, queueCapacity int) *multiJobQueue {

	if partition < 1 {
		panic("partition must bigger than 0")
	}

	multi := &multiJobQueue{
		parition: partition,
		pushIdx:  newIDGenerator(),
		popIdx:   newIDGenerator(),
	}

	for i := 0; i < partition; i++ {
		multi.queues = append(multi.queues, createJobQueue(queueCapacity))
	}
	return multi
}

// 按照 Round-Robin的方式选择队列，并放入任务
func (mjq *multiJobQueue) Push(jb JobInterface) {
	mjq.queues[mjq.pushIdx.Next()%uint64(mjq.parition)].Push(jb)
}

func (mjq *multiJobQueue) Pop() JobInterface {
	return mjq.queues[mjq.popIdx.Next()%uint64(mjq.parition)].Pop()
}

// 消费也是按照 Round-Robin的方式消费队列
func (mjq *multiJobQueue) PopTimeout(duration time.Duration) JobInterface {

	return mjq.queues[mjq.popIdx.Next()%uint64(mjq.parition)].PopTimeout(duration)
}

type idGenerator struct {
	counter uint64
}

func newIDGenerator() *idGenerator {
	return &idGenerator{
		counter: 0,
	}
}

func (g *idGenerator) Next() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}
