package user

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
	"yj-app/app/model"
	"yj-app/app/model/monitor/online"
	userModel "yj-app/app/model/system/user"
	"yj-app/app/model/system/user_post"
	"yj-app/app/model/system/user_role"
	"yj-app/app/service/middleware/sessions"
	"yj-app/app/yjgframe/cache"
	"yj-app/app/yjgframe/db"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/gconv"
	"yj-app/app/yjgframe/utils/gmd5"
	"yj-app/app/yjgframe/utils/page"
	"yj-app/app/yjgframe/utils/random"
)

//用户session列表
var SessionList sync.Map

//根据主键查询用户信息
func SelectRecordById(id int64) (*userModel.Entity, error) {
	entity := &userModel.Entity{UserId: id}
	_, err := entity.FindOne()
	return entity, err
}

// 根据条件分页查询用户列表
func SelectRecordList(param *userModel.SelectPageReq) ([]userModel.UserListEntity, *page.Paging, error) {
	return userModel.SelectPageList(param)
}

// 导出excel
func Export(param *userModel.SelectPageReq) (string, error) {
	head := []string{"用户名", "呢称", "Email", "电话号码", "性别", "部门", "领导", "状态", "删除标记", "创建人", "创建时间", "备注"}
	col := []string{"u.login_name", "u.user_name", "u.email", "u.phonenumber", "u.sex", "d.dept_name", "d.leader", "u.status", "u.del_flag", "u.create_by", "u.create_time", "u.remark"}
	return userModel.SelectExportList(param, head, col)
}

//新增用户
func AddSave(req *userModel.AddReq, c *gin.Context) (int64, error) {
	var user userModel.Entity
	user.LoginName = req.LoginName
	user.UserName = req.UserName
	user.Email = req.Email
	user.Phonenumber = req.Phonenumber
	user.Status = req.Status
	user.Sex = req.Sex
	user.DeptId = req.DeptId
	user.Remark = req.Remark

	//生成密码
	newSalt := random.GenerateSubId(6)
	newToken := req.LoginName + req.Password + newSalt
	newToken = gmd5.MustEncryptString(newToken)

	user.Salt = newSalt
	user.Password = newToken

	user.CreateTime = time.Now()

	createUser := GetProfile(c)

	if createUser != nil {
		user.CreateBy = createUser.LoginName
	}

	user.DelFlag = "0"

	session := db.Instance().Engine().NewSession()
	err := session.Begin()

	_, err = session.Table(userModel.TableName()).Insert(&user)

	if err != nil || user.UserId <= 0 {
		session.Rollback()
		return 0, err
	}

	//增加岗位数据
	if req.PostIds != "" {
		postIds := convert.ToInt64Array(req.PostIds, ",")
		userPosts := make([]user_post.Entity, 0)
		for i := range postIds {
			if postIds[i] > 0 {
				var userPost user_post.Entity
				userPost.UserId = user.UserId
				userPost.PostId = postIds[i]
				userPosts = append(userPosts, userPost)
			}
		}
		if len(userPosts) > 0 {
			_, err := session.Table(user_post.TableName()).Insert(userPosts)
			if err != nil {
				session.Rollback()
				return 0, err
			}
		}

	}

	//增加角色数据
	if req.RoleIds != "" {
		roleIds := convert.ToInt64Array(req.RoleIds, ",")
		userRoles := make([]user_role.Entity, 0)
		for i := range roleIds {
			if roleIds[i] > 0 {
				var userRole user_role.Entity
				userRole.UserId = user.UserId
				userRole.RoleId = roleIds[i]
				userRoles = append(userRoles, userRole)
			}
		}
		if len(userRoles) > 0 {
			_, err := session.Table(user_role.TableName()).Insert(userRoles)
			if err != nil {
				session.Rollback()
				return 0, err
			}
		}
	}

	return user.UserId, session.Commit()
}

//新增用户
func EditSave(req *userModel.EditReq, c *gin.Context) (int64, error) {
	user := &userModel.Entity{UserId: req.UserId}
	ok, err := user.FindOne()
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("数据不存在")
	}

	user.UserName = req.UserName
	user.Email = req.Email
	user.Phonenumber = req.Phonenumber
	user.Status = req.Status
	user.Sex = req.Sex
	user.DeptId = req.DeptId
	user.Remark = req.Remark

	user.UpdateTime = time.Now()

	updateUser := GetProfile(c)

	if updateUser != nil {
		user.UpdateBy = updateUser.LoginName
	}

	session := db.Instance().Engine().NewSession()
	tanErr := session.Begin()

	_, tanErr = session.Table(userModel.TableName()).ID(user.UserId).Update(user)

	if tanErr != nil {
		session.Rollback()
		return 0, tanErr
	}

	//增加岗位数据
	if req.PostIds != "" {
		postIds := convert.ToInt64Array(req.PostIds, ",")
		userPosts := make([]user_post.Entity, 0)
		for i := range postIds {
			if postIds[i] > 0 {
				var userPost user_post.Entity
				userPost.UserId = user.UserId
				userPost.PostId = postIds[i]
				userPosts = append(userPosts, userPost)
			}
		}
		if len(userPosts) > 0 {
			session.Exec("delete from sys_user_post where user_id=?", user.UserId)
			_, tanErr = session.Table(user_post.TableName()).Insert(userPosts)
			if tanErr != nil {
				session.Rollback()
				return 0, err
			}
		}

	}

	//增加角色数据
	if req.RoleIds != "" {
		roleIds := convert.ToInt64Array(req.RoleIds, ",")
		userRoles := make([]user_role.Entity, 0)
		for i := range roleIds {
			if roleIds[i] > 0 {
				var userRole user_role.Entity
				userRole.UserId = user.UserId
				userRole.RoleId = roleIds[i]
				userRoles = append(userRoles, userRole)
			}
		}
		if len(userRoles) > 0 {
			session.Exec("delete from sys_user_role where user_id=?", user.UserId)
			_, err := session.Table(user_role.TableName()).Insert(userRoles)
			if tanErr != nil {
				session.Rollback()
				return 0, err
			}
		}
	}

	return 1, session.Commit()
}

//根据主键删除用户信息
func DeleteRecordById(id int64) bool {
	entity := &userModel.Entity{UserId: id}
	result, _ := entity.Delete()
	if result > 0 {
		return true
	}
	return false
}

//批量删除用户记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, _ := userModel.DeleteBatch(idarr...)
	user_role.DeleteBatch(idarr...)
	user_post.DeleteBatch(idarr...)
	return result
}

//判断是否是系统管理员
func IsAdmin(userId int64) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

// 判断用户是否已经登录
func IsSignedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	uid := session.Get(model.USER_ID)
	if uid != nil {
		return true
	}
	return false
}

// 用户登录，成功返回用户信息，否则返回nil; passport应当会md5值字符串
func SignIn(loginnName, password string, c *gin.Context) (string, error) {
	//查询用户信息
	user := userModel.Entity{LoginName: loginnName}
	ok, err := user.FindOne()

	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("用户名或者密码错误")
	}

	//校验密码
	token := user.LoginName + password + user.Salt

	token = gmd5.MustEncryptString(token)

	if strings.Compare(user.Password, token) == -1 {
		return "", errors.New("密码错误")
	}
	return SaveUserToSession(&user, c), nil
}

//保存用户信息到session
func SaveUserToSession(user *userModel.Entity, c *gin.Context) string {
	session := sessions.Default(c)
	session.Set(model.USER_ID, user.UserId)
	tmp, _ := json.Marshal(user)
	session.Set(model.USER_SESSION_MARK, string(tmp))
	session.Save()
	sessionId := session.SessionId()
	SessionList.Store(sessionId, c)
	return sessionId
}

//清空用户菜单缓存
func ClearMenuCache(user *userModel.Entity) {
	if IsAdmin(user.UserId) {
		cache.Instance().Delete(model.MENU_CACHE)
	} else {
		cache.Instance().Delete(model.MENU_CACHE + gconv.String(user.UserId))
	}
}

// 用户注销
func SignOut(c *gin.Context) error {
	user := GetProfile(c)
	if user != nil {
		ClearMenuCache(user)
	}
	session := sessions.Default(c)
	sessionId := session.SessionId()

	SessionList.Delete(sessionId)
	entity := online.Entity{Sessionid: sessionId}
	entity.Delete()

	session.Delete(model.USER_ID)
	session.Delete(model.USER_SESSION_MARK)
	return session.Save()
}

//强退用户
func ForceLogout(sessionId string) error {
	var tmp interface{}
	if v, ok := SessionList.Load(sessionId); ok {
		tmp = v
	}

	if tmp != nil {
		c := tmp.(*gin.Context)
		if c != nil {
			SignOut(c)
			SessionList.Delete(sessionId)
			entity := online.Entity{Sessionid: sessionId}
			entity.Delete()
		}
	}

	return nil
}

// 检查账号是否符合规范,存在返回false,否则true
func CheckPassport(loginName string) bool {
	entity := userModel.Entity{LoginName: loginName}
	if ok, err := entity.FindOne(); err != nil || !ok {
		return false
	} else {
		return true
	}
}

// 检查登陆名是否存在,存在返回true,否则false
func CheckNickName(userName string) bool {
	entity := userModel.Entity{UserName: userName}
	if ok, err := entity.FindOne(); err != nil || !ok {
		return false
	} else {
		return true
	}
}

// 检查登陆名是否存在,存在返回true,否则false
func CheckLoginName(loginName string) bool {
	entity := userModel.Entity{LoginName: loginName}
	if ok, err := entity.FindOne(); err != nil || !ok {
		return false
	} else {
		return true
	}
}

// 获得用户信息详情
func GetProfile(c *gin.Context) *userModel.Entity {
	session := sessions.Default(c)
	tmp := session.Get(model.USER_SESSION_MARK)
	s := tmp.(string)
	var u userModel.Entity
	err := json.Unmarshal([]byte(s), &u)
	if err != nil {
		return nil
	}
	return &u
}

//更新用户信息详情
func UpdateProfile(profile *userModel.ProfileReq, c *gin.Context) error {
	user := GetProfile(c)

	if profile.UserName != "" {
		user.UserName = profile.UserName
	}

	if profile.Email != "" {
		user.Email = profile.Email
	}

	if profile.Phonenumber != "" {
		user.Phonenumber = profile.Phonenumber
	}

	if profile.Sex != "" {
		user.Sex = profile.Sex
	}

	_, err := user.Update()
	if err != nil {
		return errors.New("保存数据失败")
	}

	SaveUserToSession(user, c)
	return nil
}

//更新用户头像
func UpdateAvatar(avatar string, c *gin.Context) error {
	user := GetProfile(c)

	if avatar != "" {
		user.Avatar = avatar
	}

	_, err := user.Update()
	if err != nil {
		return errors.New("保存数据失败")
	}

	SaveUserToSession(user, c)
	return nil
}

//修改用户密码
func UpdatePassword(profile *userModel.PasswordReq, c *gin.Context) error {
	user := GetProfile(c)

	if profile.OldPassword == "" {
		return errors.New("旧密码不能为空")
	}

	if profile.NewPassword == "" {
		return errors.New("新密码不能为空")
	}

	if profile.Confirm == "" {
		return errors.New("确认密码不能为空")
	}

	if profile.NewPassword == profile.OldPassword {
		return errors.New("新旧密码不能相同")
	}

	if profile.Confirm != profile.NewPassword {
		return errors.New("确认密码不一致")
	}

	//校验密码
	token := user.LoginName + profile.OldPassword + user.Salt
	token = gmd5.MustEncryptString(token)

	if token != user.Password {
		return errors.New("原密码不正确")
	}

	//新校验密码
	newSalt := random.GenerateSubId(6)
	newToken := user.LoginName + profile.NewPassword + newSalt
	newToken = gmd5.MustEncryptString(newToken)

	user.Salt = newSalt
	user.Password = newToken

	_, err := user.Update()
	if err != nil {
		return errors.New("保存数据失败")
	}

	SaveUserToSession(user, c)
	return nil
}

//重置用户密码
func ResetPassword(params *userModel.ResetPwdReq) (bool, error) {
	user := userModel.Entity{UserId: params.UserId}
	if ok, err := user.FindOne(); err != nil || !ok {
		return false, errors.New("用户不存在")
	}
	//新校验密码
	newSalt := random.GenerateSubId(6)
	newToken := user.LoginName + params.Password + newSalt
	newToken = gmd5.MustEncryptString(newToken)

	user.Salt = newSalt
	user.Password = newToken
	if _, err := user.Update(); err != nil {
		return false, errors.New("保存数据失败")
	}
	return true, nil
}

//校验密码是否正确
func CheckPassword(user *userModel.Entity, password string) bool {
	if user == nil || user.UserId <= 0 {
		return false
	}

	//校验密码
	token := user.LoginName + password + user.Salt
	token = gmd5.MustEncryptString(token)

	if strings.Compare(token, user.Password) == 0 {
		return true
	} else {
		return false
	}
}

//检查邮箱是否已使用
func CheckEmailUnique(userId int64, email string) bool {
	return userModel.CheckEmailUnique(userId, email)
}

//检查邮箱是否存在,存在返回true,否则false
func CheckEmailUniqueAll(email string) bool {
	return userModel.CheckEmailUniqueAll(email)
}

//检查手机号是否已使用,存在返回true,否则false
func CheckPhoneUnique(userId int64, phone string) bool {
	return userModel.CheckPhoneUnique(userId, phone)
}

//检查手机号是否已使用 ,存在返回true,否则false
func CheckPhoneUniqueAll(phone string) bool {
	return userModel.CheckPhoneUniqueAll(phone)
}

//根据登陆名查询用户信息
func SelectUserByLoginName(loginName string) (*userModel.Entity, error) {
	return userModel.SelectUserByLoginName(loginName)
}

//根据手机号查询用户信息
func SelectUserByPhoneNumber(phonenumber string) (*userModel.Entity, error) {
	return userModel.SelectUserByPhoneNumber(phonenumber)
}

// 查询已分配用户角色列表
func SelectAllocatedList(roleId int64, loginName, phonenumber string) ([]userModel.Entity, error) {
	return userModel.SelectAllocatedList(roleId, loginName, phonenumber)
}

// 查询未分配用户角色列表
func SelectUnallocatedList(roleId int64, loginName, phonenumber string) ([]userModel.Entity, error) {
	return userModel.SelectUnallocatedList(roleId, loginName, phonenumber)
}
