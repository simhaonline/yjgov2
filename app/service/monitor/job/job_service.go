// ==========================================================================
// 云捷GO自动生成业务逻辑层相关代码，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2020-02-18 15:44:13
// 生成路径: app/service/module/job/job_service.go
// 生成人：yunjie
// ==========================================================================
package job

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	jobModel "yj-app/app/model/monitor/job"
	userService "yj-app/app/service/system/user"
	"yj-app/app/task"
	"yj-app/app/yjgframe/cron"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/gconv"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*jobModel.Entity, error) {
	entity := &jobModel.Entity{JobId: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	entity := &jobModel.Entity{JobId: id}
	result, err := entity.Delete()
	if err == nil && result > 0 {
		return true
	}

	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, _ := jobModel.DeleteBatch(idarr...)
	return result
}

//添加数据
func AddSave(req *jobModel.AddReq, c *gin.Context) (int64, error) {
	//检查任务名称是否存在
	rs := cron.Get(req.JobName)

	if rs != nil {
		return 0, errors.New("任务名称已经存在")
	}

	//可以task目录下是否绑定对应的方法
	f := task.GetByName(req.JobName)
	if f == nil {
		return 0, errors.New("当前task目录下没有绑定这个方法")
	}

	var entity jobModel.Entity
	entity.JobName = req.JobName
	entity.JobParams = req.JobParams
	entity.JobGroup = req.JobGroup
	entity.InvokeTarget = req.InvokeTarget
	entity.CronExpression = req.CronExpression
	entity.MisfirePolicy = req.MisfirePolicy
	entity.Concurrent = req.Concurrent
	entity.Status = req.Status
	entity.Remark = req.Remark
	entity.CreateTime = time.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(c)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	entity.Insert()
	return entity.JobId, nil
}

//修改数据
func EditSave(req *jobModel.EditReq, c *gin.Context) (int64, error) {
	//检查任务名称是否存在
	tmp := cron.Get(req.JobName)

	if tmp != nil {
		tmp.Stop()
	}

	//可以task目录下是否绑定对应的方法
	f := task.GetByName(req.JobName)
	if f == nil {
		return 0, errors.New("当前task目录下没有绑定这个方法")
	}

	entity := &jobModel.Entity{JobId: req.JobId}
	_, err := entity.FindOne()

	if err != nil {
		return 0, err
	}

	if entity == nil {
		return 0, errors.New("数据不存在")
	}

	entity.InvokeTarget = req.InvokeTarget
	entity.JobParams = req.JobParams
	entity.CronExpression = req.CronExpression
	entity.MisfirePolicy = req.MisfirePolicy
	entity.Concurrent = req.Concurrent
	entity.Status = req.Status
	entity.Remark = req.Remark
	entity.UpdateTime = time.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(c)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	return entity.Update()
}

//初始化任务状态
func Init() {
	list, err := jobModel.SelectListAll(nil)
	if err != nil {
		return
	}

	stopIds := ""
	startIds := ""

	for i := 0; i < len(list); i++ {
		if len(list[i].JobName) > 0 {
			rs := cron.Get(list[i].JobName)
			if list[i].Status == "0" && rs == nil {
				if stopIds == "" {
					stopIds = gconv.String(list[i].JobId)
				} else {
					stopIds += "," + gconv.String(list[i].JobId)
				}
			}

			if list[i].Status == "1" && rs != nil {
				if startIds == "" {
					startIds = gconv.String(list[i].JobId)
				} else {
					startIds += "," + gconv.String(list[i].JobId)
				}
			}
		}
	}

	if stopIds != "" {
		jobModel.UpdateState(stopIds, "1")
	}

	if startIds != "" {
		jobModel.UpdateState(startIds, "0")
	}

}

//启动任务
func Start(entity *jobModel.Entity) error {
	//可以task目录下是否绑定对应的方法
	f := task.GetByName(entity.JobName)
	if f == nil {
		return errors.New("当前task目录下没有绑定这个方法")
	}

	//传参
	paramArr := strings.Split(entity.JobParams, "|")
	task.EditParams(f.FuncName, paramArr)

	rs := cron.Get(entity.JobName)

	if rs == nil {
		if entity.MisfirePolicy == "1" {
			j, err := cron.New(entity, f.Run)

			if err != nil && j == nil {
				return err
			}

			entity.Status = "0"
			entity.Update()
		} else {
			f.Run()
		}
	} else {
		return errors.New("任务已存在")
	}

	return nil
}

//停止任务
func Stop(entity *jobModel.Entity) error {
	//可以task目录下是否绑定对应的方法
	f := task.GetByName(entity.JobName)
	if f == nil {
		return errors.New("当前task目录下没有绑定这个方法")
	}

	rs := cron.Get(entity.JobName)

	if rs != nil {
		rs.Stop()
	}

	entity.Status = "1"
	entity.Update()
	return nil
}

//根据条件查询数据
func SelectListAll(params *jobModel.SelectPageReq) ([]jobModel.Entity, error) {
	return jobModel.SelectListAll(params)
}

//根据条件分页查询数据
func SelectListByPage(params *jobModel.SelectPageReq) (*[]jobModel.Entity, *page.Paging, error) {
	return jobModel.SelectListByPage(params)
}

// 导出excel
func Export(param *jobModel.SelectPageReq) (string, error) {

	head := []string{"任务ID", "任务名称", "任务组名", "调用目标字符串", "cron执行表达式", "计划执行错误策略（1立即执行 2执行一次 3放弃执行）", "是否并发执行（0允许 1禁止）", "状态（0正常 1暂停）", "创建者", "创建时间", "更新者", "更新时间", "备注信息"}
	col := []string{"job_id", "job_name", "job_group", "invoke_target", "cron_expression", "misfire_policy", "concurrent", "status", "create_by", "create_time", "update_by", "update_time", "remark"}
	return jobModel.SelectListExport(param, head, col)
}
