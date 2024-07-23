package easyqueue

import "fmt"

var ErrOverFlow error = fmt.Errorf("queue is overflow")

type WaitJob interface {
	Done() <-chan struct{}
	Err() error
}

type JobInterface interface {
	WaitJob
	DoJob()
}

// 创建任务
func newJob(fn func()) JobInterface {

	return &job{
		fn:  fn,
		ch:  make(chan struct{}),
		err: nil,
	}

}

type job struct {
	fn  func()
	err error
	ch  chan struct{}
}

func (h *job) DoJob() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		close(h.ch)
	}()

	h.fn()
}

func (h *job) Err() error {
	return h.err
}

func (h *job) Done() <-chan struct{} {
	return h.ch
}
