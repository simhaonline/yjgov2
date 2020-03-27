package dept

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yj-app/app/model"
	deptModel "yj-app/app/model/system/dept"
	deptService "yj-app/app/service/system/dept"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/utils/gconv"
)

//列表页
func List(c *gin.Context) {
	response.BuildTpl(c, "system/dept/list").WriteTpl()
}

//列表分页数据
func ListAjax(c *gin.Context) {
	var req *deptModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("部门管理", req).WriteJsonExit()
		return
	}
	rows := make([]deptModel.Entity, 0)
	result, err := deptService.SelectListAll(req)

	if err == nil && len(result) > 0 {
		rows = result
	}

	c.JSON(http.StatusOK, rows)
}

//新增页面
func Add(c *gin.Context) {
	pid := gconv.Int64(c.Query("pid"))

	if pid == 0 {
		pid = 100
	}

	tmp := deptService.SelectDeptById(pid)

	response.BuildTpl(c, "system/dept/add").WriteTpl(gin.H{"dept": tmp})
}

//新增页面保存
func AddSave(c *gin.Context) {
	var req *deptModel.AddReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg(err.Error()).Log("部门管理", req).WriteJsonExit()
		return
	}

	if deptService.CheckDeptNameUniqueAll(req.DeptName, req.ParentId) == "1" {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg("部门名称已存在").Log("部门管理", req).WriteJsonExit()
		return
	}

	rid, err := deptService.AddSave(req, c)

	if err != nil || rid <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).Log("部门管理", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetBtype(model.Buniss_Add).Log("部门管理", req).WriteJsonExit()
}

//修改页面
func Edit(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	dept := deptService.SelectDeptById(id)

	if dept == nil || dept.DeptId <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "部门不存在",
		})
		return
	}

	response.BuildTpl(c, "system/dept/edit").WriteTpl(gin.H{
		"dept": dept,
	})
}

//修改页面保存
func EditSave(c *gin.Context) {
	var req *deptModel.EditReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("部门管理", req).WriteJsonExit()
		return
	}

	if deptService.CheckDeptNameUnique(req.DeptName, req.DeptId, req.ParentId) == "1" {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg("部门名称已存在").Log("部门管理", req).WriteJsonExit()
		return
	}

	rs, err := deptService.EditSave(req, c)

	if err != nil || rs <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).Log("部门管理", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetData(rs).SetBtype(model.Buniss_Edit).Log("部门管理", req).WriteJsonExit()
}

//删除数据
func Remove(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))
	rs := deptService.DeleteDeptById(id)

	if rs > 0 {
		response.SucessResp(c).SetBtype(model.Buniss_Del).Log("部门管理", gin.H{"id": id}).WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("部门管理", gin.H{"id": id}).WriteJsonExit()
	}
}

//加载部门列表树结构的数据
func TreeData(c *gin.Context) {
	result, _ := deptService.SelectDeptTree(0, "", "")
	c.JSON(http.StatusOK, result)
}

//加载部门列表树选择页面
func SelectDeptTree(c *gin.Context) {
	deptId := gconv.Int64(c.Query("deptId"))
	deptPoint := deptService.SelectDeptById(deptId)

	if deptPoint != nil {
		response.BuildTpl(c, "system/dept/tree").WriteTpl(gin.H{
			"dept": *deptPoint,
		})
	} else {
		response.BuildTpl(c, "system/dept/tree").WriteTpl()
	}
}

//加载角色部门（数据权限）列表树
func RoleDeptTreeData(c *gin.Context) {
	roleId := gconv.Int64(c.Query("roleId"))
	result, err := deptService.RoleDeptTreeData(roleId)

	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("菜单树", gin.H{"roleId": roleId})
		return
	}

	c.JSON(http.StatusOK, result)
}

//检查部门名称是否已经存在
func CheckDeptNameUnique(c *gin.Context) {
	var req *deptModel.CheckDeptNameReq
	if err := c.ShouldBind(&req); err != nil {
		c.Writer.WriteString("1")
		return
	}

	result := deptService.CheckDeptNameUnique(req.DeptName, req.DeptId, req.ParentId)

	c.Writer.WriteString(result)
}

//检查部门名称是否已经存在
func CheckDeptNameUniqueAll(c *gin.Context) {
	var req *deptModel.CheckDeptNameALLReq
	if err := c.ShouldBind(&req); err != nil {
		c.Writer.WriteString("1")
		return
	}

	result := deptService.CheckDeptNameUniqueAll(req.DeptName, req.ParentId)

	c.Writer.WriteString(result)
}
