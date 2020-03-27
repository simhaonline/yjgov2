package logininfor

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/model"
	logininforModel "yj-app/app/model/monitor/logininfor"
	logininforService "yj-app/app/service/monitor/logininfor"
	"yj-app/app/yjgframe/response"
)

//用户列表页
func List(c *gin.Context) {
	response.BuildTpl(c, "monitor/logininfor/list").WriteTpl()
}

//用户列表分页数据
func ListAjax(c *gin.Context) {
	var req *logininforModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}

	rows := make([]logininforModel.Entity, 0)

	result, page, err := logininforService.SelectPageList(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}
	response.BuildTable(c, page.Total, rows).WriteJsonExit()
}

//删除数据
func Remove(c *gin.Context) {
	var req *model.RemoveReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).SetMsg(err.Error()).Log("登陆日志管理", req).WriteJsonExit()
		return
	}

	rs := logininforService.DeleteRecordByIds(req.Ids)

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).SetData(rs).Log("登陆日志管理", req).WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("登陆日志管理", req).WriteJsonExit()
	}
}

//清空记录
func Clean(c *gin.Context) {

	rs, _ := logininforService.DeleteRecordAll()

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).SetData(rs).Log("登陆日志管理", "all").WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("登陆日志管理", "all").WriteJsonExit()
	}
}

//导出
func Export(c *gin.Context) {
	var req *logininforModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出登陆日志", req).WriteJsonExit()
		return
	}

	url, err := logininforService.Export(req)

	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出登陆日志", req).WriteJsonExit()
	} else {
		response.SucessResp(c).SetMsg(url).Log("导出登陆日志", req).WriteJsonExit()
	}
}

//解锁账号
func Unlock(c *gin.Context) {
	loginName := c.Query("loginName")
	if loginName == "" {
		response.ErrorResp(c).SetMsg("参数错误").Log("解锁账号", "loginName="+loginName).WriteJsonExit()
	} else {
		logininforService.RemovePasswordCounts(loginName)
		logininforService.Unlock(loginName)
		response.SucessResp(c).Log("解锁账号", "loginName="+loginName).WriteJsonExit()
	}

}
