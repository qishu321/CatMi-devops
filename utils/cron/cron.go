package cron

import (
	"time"

	"github.com/go-co-op/gocron"
)

type JobFunc func() error
type Cron struct {
	limitDailyJobs map[string][]JobFunc
	absTimeJobs    map[int64][]JobFunc
	Scheduler      *gocron.Scheduler
	NotifyCh       chan []JobFunc
}

type Job struct {
	AbsTime time.Time
}

func NewCron() *Cron {
	sch := gocron.NewScheduler(time.Local)
	return &Cron{
		limitDailyJobs: make(map[string][]JobFunc, 3),
		absTimeJobs:    make(map[int64][]JobFunc, 3),
		Scheduler:      sch,
		NotifyCh:       make(chan []JobFunc, 10)}
}

// 绝对时间定时器任务

func (cj *Cron) AddAbsTimeJob(startTime time.Time, jFunc JobFunc) {
	ts := startTime.Unix()
	if _, has := cj.absTimeJobs[ts]; !has {
		cj.Scheduler.Every(1).LimitRunsTo(1).StartAt(time.Unix(ts, 0)).Do(func() {
			cj.NotifyCh <- cj.absTimeJobs[ts]
		})
	}

	cj.absTimeJobs[ts] = append(cj.absTimeJobs[ts], jFunc)
}

// 有限制天数的每天固定时间定时器任务

func (cj *Cron) AddLimitDailyJob(startTime time.Time, repeat int, jFunc JobFunc) {
	if startTime.Before(time.Now()) { //开始时间比现在早 需计算剩余执行次数
		daysBefore := int(time.Now().Sub(startTime).Hours()) / 24
		if int(time.Now().Sub(startTime).Hours())%24 > 0 {
			daysBefore += 1
		}

		repeat -= daysBefore
	}

	if repeat <= 0 {
		return
	}

	timeOnly := startTime.Format(time.TimeOnly)
	if _, has := cj.limitDailyJobs[timeOnly]; !has {
		cj.Scheduler.Every(1).Day().LimitRunsTo(repeat).At(timeOnly).Do(func() {
			cj.NotifyCh <- cj.limitDailyJobs[timeOnly]
		})
	}

	cj.limitDailyJobs[timeOnly] = append(cj.limitDailyJobs[timeOnly], jFunc)
}

var CronJob *Cron

func init() {
	CronJob = NewCron()
	CronJob.Scheduler.StartAsync()
	go handler()
}

func handler() {
	for {
		select {
		case jobFuncS := <-CronJob.NotifyCh:
			for _, fun := range jobFuncS {
				_ = fun()
			}
		}
	}
}
