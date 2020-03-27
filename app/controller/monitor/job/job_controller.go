package job

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yj-app/app/model"
	jobModel "yj-app/app/model/monitor/job"
	jobLogModel "yj-app/app/model/monitor/job_log"
	jobService "yj-app/app/service/monitor/job"
	jobLogService "yj-app/app/service/monitor/job_log"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/gconv"
)

//列表页
func List(c *gin.Context) {
	jobService.Init()
	response.BuildTpl(c, "monitor/job/list").WriteTpl()
}

//列表分页数据
func ListAjax(c *gin.Context) {
	var req *jobModel.SelectPageReq
	//获取参数

	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}
	rows := make([]jobModel.Entity, 0)
	result, page, err := jobService.SelectListByPage(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}
	response.BuildTable(c, page.Total, rows).WriteJsonExit()
}

//列表页
func LogList(c *gin.Context) {
	response.BuildTpl(c, "monitor/job/jobLog").WriteTpl()
}

//列表分页数据
func LogListAjax(c *gin.Context) {
	var req *jobLogModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("任务日志管理", req).WriteJsonExit()
		return
	}
	rows := make([]jobLogModel.Entity, 0)
	result, page, err := jobLogService.SelectListByPage(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}

	c.JSON(http.StatusOK, model.TableDataInfo{
		Code:  0,
		Msg:   "操作成功",
		Total: page.Total,
		Rows:  rows,
	})
}

//新增页面
func Add(c *gin.Context) {
	user := userService.GetProfile(c)
	response.BuildTpl(c, "monitor/job/add").WriteTpl(gin.H{"loginName": user.LoginName})
}

//新增页面保存
func AddSave(c *gin.Context) {
	var req *jobModel.AddReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}

	id, err := jobService.AddSave(req, c)

	if err != nil || id <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).Log("定时任务管理", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetBtype(model.Buniss_Add).SetData(id).Log("定时任务管理", req).WriteJsonExit()
}

//修改页
func Edit(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	entity, err := jobService.SelectRecordById(id)

	if err != nil || entity == nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "数据不存在",
		})
		return
	}

	user := userService.GetProfile(c)

	response.BuildTpl(c, "monitor/job/edit").WriteTpl(gin.H{
		"job":       entity,
		"loginName": user.LoginName,
	})
}

//修改页面保存
func EditSave(c *gin.Context) {
	var req jobModel.EditReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}

	rs, err := jobService.EditSave(&req, c)

	if err != nil || rs <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).Log("定时任务管理", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetBtype(model.Buniss_Edit).SetData(rs).Log("定时任务管理", req).WriteJsonExit()
}

//详情页
func Detail(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}
	job, err := jobService.SelectRecordById(id)
	if err != nil || job == nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "数据不存在",
		})
		return
	}
	response.BuildTpl(c, "monitor/job/detail").WriteTpl(gin.H{"job": job})
}

//详情页
func DetailLog(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}
	jobLog, err := jobLogService.SelectRecordById(id)
	if err != nil || jobLog == nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "数据不存在",
		})
		return
	}
	response.BuildTpl(c, "monitor/job/detailLog").WriteTpl(gin.H{"jobLog": jobLog})
}

//删除数据
func Remove(c *gin.Context) {
	var req *model.RemoveReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}

	idarr := convert.ToInt64Array(req.Ids, ",")
	list, _ := jobModel.FindIn("job_id", idarr)
	if list != nil && len(list) > 0 {
		//for _, j := range *list {
		//	gcron.Remove(j.JobName)
		//}
	}

	rs := jobService.DeleteRecordByIds(req.Ids)

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).SetData(rs).Log("定时任务管理", req).WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("定时任务管理", req).WriteJsonExit()
	}
}

//导出
func Export(c *gin.Context) {
	var req *jobModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}
	url, err := jobService.Export(req)

	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("定时任务管理", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetMsg(url).Log("定时任务管理", req).WriteJsonExit()
}

//启动
func Start(c *gin.Context) {
	jobId := gconv.Int64(c.PostForm("jobId"))
	if jobId <= 0 {
		response.ErrorResp(c).SetMsg("参数错误").Log("定时任务管理启动", gin.H{"jobId": jobId}).WriteJsonExit()
		return
	}
	job, _ := jobService.SelectRecordById(jobId)
	if job == nil {
		response.ErrorResp(c).SetMsg("任务不存在").Log("定时任务管理启动", gin.H{"jobId": jobId}).WriteJsonExit()
		return
	}

	err := jobService.Start(job)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("定时任务管理启动", gin.H{"jobId": jobId}).WriteJsonExit()
	} else {
		response.SucessResp(c).Log("定时任务管理启动", gin.H{"jobId": jobId}).WriteJsonExit()
	}
}

//停止
func Stop(c *gin.Context) {
	jobId := gconv.Int64(c.PostForm("jobId"))
	if jobId <= 0 {
		response.ErrorResp(c).SetMsg("参数错误").Log("定时任务管理停止", gin.H{"jobId": jobId}).WriteJsonExit()
		return
	}
	job, _ := jobService.SelectRecordById(jobId)
	if job == nil {
		response.ErrorResp(c).SetMsg("任务不存在").Log("定时任务管理停止", gin.H{"jobId": jobId}).WriteJsonExit()
		return
	}

	err := jobService.Stop(job)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("定时任务管理停止", gin.H{"jobId": jobId}).WriteJsonExit()
	} else {
		response.SucessResp(c).SetMsg("停止成功").Log("定时任务管理停止", gin.H{"jobId": jobId}).WriteJsonExit()
	}
}
