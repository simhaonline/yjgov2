package operlog

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"yj-app/app/model"
	operlogModel "yj-app/app/model/monitor/oper_log"
	operlogService "yj-app/app/service/monitor/operlog"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/utils/gconv"
)

//用户列表页
func List(c *gin.Context) {
	response.BuildTpl(c, "monitor/operlog/list").WriteTpl()
}

//用户列表分页数据
func ListAjax(c *gin.Context) {
	var req *operlogModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.CommonRes{
			Code: 500,
			Msg:  err.Error(),
		})
	}

	rows := make([]operlogModel.Entity, 0)

	result, page, err := operlogService.SelectPageList(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}

	response.BuildTable(c, page.Total, rows).WriteJsonExit()
}

//清空记录
func Clean(c *gin.Context) {

	rs, _ := operlogService.DeleteRecordAll()

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).SetData(rs).Log("操作日志管理", "all").WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("操作日志管理", "all").WriteJsonExit()
	}
}

//删除数据
func Remove(c *gin.Context) {
	var req *model.RemoveReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).SetMsg(err.Error()).Log("操作日志管理", req).WriteJsonExit()
		return
	}

	rs := operlogService.DeleteRecordByIds(req.Ids)

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).SetData(rs).Log("操作日志管理", req).WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("操作日志管理", req).WriteJsonExit()
	}
}

//记录详情
func Detail(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	operLog, err := operlogService.SelectRecordById(id)

	if err != nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "数据不存在",
		})
		return
	}

	jsonResult := template.HTML(operLog.JsonResult)
	operParam := template.HTML(operLog.OperParam)
	response.BuildTpl(c, "monitor/operlog/detail").WriteTpl(gin.H{
		"operLog":    operLog,
		"jsonResult": jsonResult,
		"operParam":  operParam,
	})
}

//导出
func Export(c *gin.Context) {
	var req *operlogModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出操作日志", req).WriteJsonExit()
		return
	}
	url, err := operlogService.Export(req)

	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出操作日志", req).WriteJsonExit()
	} else {
		response.SucessResp(c).SetMsg(url).Log("导出操作日志", req).WriteJsonExit()
	}
}
