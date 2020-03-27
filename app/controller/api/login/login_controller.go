package login

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/token"
	"yj-app/app/yjgframe/utils/gconv"
)

type user struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// @Summary 登陆
// @Description api测试
// @Accept  json
// @Produce  json
// @Success 200 {object} model.CommonRes
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	u := new(user)
	if err := c.ShouldBind(&u); err != nil {
		response.ErrorResp(c).SetData(err.Error()).WriteJsonExit()
		return
	}

	//验证用户名和密码
	if u == nil || u.Username == "" || u.Password == "" {
		response.ErrorResp(c).SetData("用户名或密码不正确").WriteJsonExit()
		return
	}

	//获取用户id
	uid := 10

	//生成token
	token, err := token.New(gconv.String(uid)).CreateToken()

	if err != nil {
		response.ErrorResp(c).SetData("Error while signing the token").WriteJsonExit()
		return
	}

	//返回token
	response.SucessResp(c).SetData(token).WriteJsonExit()
}

// @Summary api测试
// @Description api测试
// @Accept  json
// @Produce  json
// @Success 200 {object} model.CommonRes
// @Router /api/v1/loginOut [get]
func Test(c *gin.Context) {
	uid, _ := c.Get("uid")
	response.SucessResp(c).SetData(uid).WriteJsonExit()
}
