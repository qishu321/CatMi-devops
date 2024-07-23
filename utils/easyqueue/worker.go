package easyqueue

import (
	"sync"
	"sync/atomic"
	"time"
)

type worker struct {
	queue  Queue        // 工作协程监视的队列
	closed atomic.Int32 // 是否停止工作协程标识
}

func createWorker(queue Queue) *worker {
	return &worker{
		queue:  queue,
		closed: atomic.Int32{},
	}
}

func (w *worker) Run() {

	w.closed.Store(0)
	go func() {
		for {

			if w.closed.Load() > 0 { // 关闭worker
				break
			}

			job := w.queue.PopTimeout(1 * time.Millisecond)
			//job := w.queue.Pop()
			if job != nil {
				job.DoJob()
			}
		}
	}()
}

func (w *worker) Stop() {
	w.closed.Store(1)
}

// 消费组
type workerGroup struct {
	mu         sync.Mutex
	workers    []*worker
	workersNum int
	queue      Queue
}

func createWorkerMange(queue Queue, workerNum int) *workerGroup {
	mange := workerGroup{
		workersNum: workerNum,
		queue:      queue,
	}
	for i := 0; i < workerNum; i++ {
		worker := createWorker(queue)
		worker.Run()
		mange.workers = append(mange.workers, worker)
	}
	return &mange
}

// 动态调整消费者的数量
func (wm *workerGroup) Adjust(workerNum int) {
	if workerNum == wm.workersNum {
		return // do nothing
	}

	wm.mu.Lock()
	defer wm.mu.Unlock()

	if workerNum > wm.workersNum { // 增加
		for i := wm.workersNum; i < workerNum; i++ {
			worker := createWorker(wm.queue)
			worker.Run()
			wm.workers = append(wm.workers, worker)
		}

	} else { // 减少
		for i := workerNum; i < wm.workersNum; i++ {
			wm.workers[i].Stop()
		}
		wm.workers = wm.workers[:wm.workersNum-workerNum]
	}

	wm.workersNum = workerNum
}
