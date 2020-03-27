package user

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/model"
	postModel "yj-app/app/model/system/post"
	roleModel "yj-app/app/model/system/role"
	userModel "yj-app/app/model/system/user"
	deptServic "yj-app/app/service/system/dept"
	postService "yj-app/app/service/system/post"
	roleService "yj-app/app/service/system/role"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/utils/gconv"
)

//用户列表页
func List(c *gin.Context) {
	response.BuildTpl(c, "system/user/list").WriteTpl()
}

//用户列表分页数据
func ListAjax(c *gin.Context) {
	var req *userModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("用户管理", req).WriteJsonExit()
		return
	}
	rows := make([]userModel.UserListEntity, 0)
	result, page, err := userService.SelectRecordList(req)

	if err == nil && len(result) > 0 {
		rows = result
	}
	response.BuildTable(c, page.Total, rows).WriteJsonExit()
}

//用户新增页面
func Add(c *gin.Context) {
	var paramsRole *roleModel.SelectPageReq
	var paramsPost *postModel.SelectPageReq

	roles := make([]roleModel.EntityFlag, 0)
	posts := make([]postModel.EntityFlag, 0)

	rolesP, _ := roleService.SelectRecordAll(paramsRole)

	if rolesP != nil {
		roles = rolesP
	}

	postP, _ := postService.SelectListAll(paramsPost)

	if postP != nil {
		posts = postP
	}
	response.BuildTpl(c, "system/user/add").WriteTpl(gin.H{
		"roles": roles,
		"posts": posts,
	})
}

//保存新增用户数据
func AddSave(c *gin.Context) {
	var req *userModel.AddReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg(err.Error()).Log("新增用户", req).WriteJsonExit()
		return
	}

	//判断登陆名是否已注册
	isHadName := userService.CheckLoginName(req.LoginName)
	if isHadName {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg("登陆名已经存在").Log("新增用户", req).WriteJsonExit()
		return
	}

	//判断手机号码是否已注册
	isHadPhone := userService.CheckPhoneUniqueAll(req.Phonenumber)
	if isHadPhone {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg("手机号码已经存在").Log("新增用户", req).WriteJsonExit()
		return
	}

	//判断邮箱是否已注册
	isHadEmail := userService.CheckEmailUniqueAll(req.Email)
	if isHadEmail {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).SetMsg("邮箱已经存在").Log("新增用户", req).WriteJsonExit()
		return
	}

	uid, err := userService.AddSave(req, c)

	if err != nil || uid <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Add).Log("新增用户", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetData(uid).SetBtype(model.Buniss_Add).Log("新增用户", req).WriteJsonExit()
}

//用户修改页面
func Edit(c *gin.Context) {
	id := gconv.Int64(c.Query("id"))

	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	user, err := userService.SelectRecordById(id)

	if err != nil || user == nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "用户不存在",
		})
		return
	}

	//获取部门信息
	deptName := ""
	if user.DeptId > 0 {
		dept := deptServic.SelectDeptById(user.DeptId)
		if dept != nil {
			deptName = dept.DeptName
		}
	}

	roles := make([]roleModel.EntityFlag, 0)
	posts := make([]postModel.EntityFlag, 0)

	rolesP, _ := roleService.SelectRoleContactVo(id)

	if rolesP != nil {
		roles = rolesP
	}

	postP, _ := postService.SelectPostsByUserId(id)

	if postP != nil {
		posts = postP
	}

	response.BuildTpl(c, "system/user/edit").WriteTpl(gin.H{
		"user":     user,
		"deptName": deptName,
		"roles":    roles,
		"posts":    posts,
	})
}

//重置密码
func ResetPwd(c *gin.Context) {
	id := gconv.Int64(c.Query("userId"))
	if id <= 0 {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	user, err := userService.SelectRecordById(id)

	if err != nil || user == nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "用户不存在",
		})
		return
	}
	response.BuildTpl(c, "system/user/resetPwd").WriteTpl(gin.H{
		"user": user,
	})
}

//重置密码保存
func ResetPwdSave(c *gin.Context) {
	var req *userModel.ResetPwdReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("重置密码", req).WriteJsonExit()
	}

	result, err := userService.ResetPassword(req)

	if err != nil || !result {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("重置密码", req).WriteJsonExit()
	} else {
		response.SucessResp(c).SetBtype(model.Buniss_Edit).Log("重置密码", req).WriteJsonExit()
	}
}

//保存修改用户数据
func EditSave(c *gin.Context) {
	var req *userModel.EditReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("修改用户", req).WriteJsonExit()
		return
	}

	//判断手机号码是否已注册
	isHadPhone := userService.CheckPhoneUnique(req.UserId, req.Phonenumber)
	if isHadPhone {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg("手机号码已经存在").Log("修改用户", req).WriteJsonExit()
		return
	}

	//判断邮箱是否已注册
	isHadEmail := userService.CheckEmailUnique(req.UserId, req.Email)
	if isHadEmail {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).SetMsg("邮箱已经存在").Log("修改用户", req).WriteJsonExit()
		return
	}

	uid, err := userService.EditSave(req, c)

	if err != nil || uid <= 0 {
		response.ErrorResp(c).SetBtype(model.Buniss_Edit).Log("修改用户", req).WriteJsonExit()
		return
	}

	response.SucessResp(c).SetData(uid).SetBtype(model.Buniss_Edit).Log("修改用户", req).WriteJsonExit()
}

//删除数据
func Remove(c *gin.Context) {
	var req *model.RemoveReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).SetMsg(err.Error()).Log("删除用户", req).WriteJsonExit()
	}

	rs := userService.DeleteRecordByIds(req.Ids)

	if rs > 0 {
		response.SucessResp(c).SetData(rs).SetBtype(model.Buniss_Del).Log("删除用户", req).WriteJsonExit()
	} else {
		response.ErrorResp(c).SetBtype(model.Buniss_Del).Log("删除用户", req).WriteJsonExit()
	}
}

//导出
func Export(c *gin.Context) {
	var req *userModel.SelectPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
	}
	url, err := userService.Export(req)

	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
		return
	}
	response.SucessResp(c).SetMsg(url).Log("导出Excel", req).WriteJsonExit()
}
