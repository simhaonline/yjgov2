package cron

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sync"
	"yj-app/app/model/monitor/job"
)

var (
	mainCron *cron.Cron
	jobList  = sync.Map{}
)

type jobStruct struct {
	JobName string
	JobId   int
	JobCron *cron.Cron
}

//新建一个任务
func New(job *job.Entity, cmd func()) (*jobStruct, error) {
	_, ok := jobList.Load(job.JobName)
	if ok {
		return nil, errors.New("任务已存在")
	}

	jobCron := cron.New()
	id, err := jobCron.AddFunc(job.CronExpression, cmd)
	if err != nil {
		return nil, err
	}
	jobCron.Start()

	var js jobStruct
	js.JobId = int(id)
	js.JobName = job.JobName
	js.JobCron = jobCron

	jobList.Store(job.JobName, &js)
	return &js, nil
}

//根据任务名称获取任务
func Get(jobName string) *jobStruct {
	j, ok := jobList.Load(jobName)
	if !ok {
		return nil
	}
	result, sucessed := j.(*jobStruct)
	if !sucessed {
		return nil
	}
	return result
}

//停止任务
func (j *jobStruct) Stop() {
	j.JobCron.Stop()
	j.JobCron.Remove(cron.EntryID(j.JobId))
	jobList.Delete(j.JobName)
}
