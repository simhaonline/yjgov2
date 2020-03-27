package index

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"yj-app/app/model"
	"yj-app/app/model/system/menu"
	configService "yj-app/app/service/system/config"
	menuService "yj-app/app/service/system/menu"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/response"
)

//后台框架首页
func Index(c *gin.Context) {
	if userService.IsSignedIn(c) {
		user := userService.GetProfile(c)
		loginname := user.LoginName
		username := user.UserName
		avatar := user.Avatar
		if avatar == "" {
			avatar = "/resource/img/profile.jpg"
		}

		var menus *[]menu.EntityExtend

		//获取菜单数据
		if userService.IsAdmin(user.UserId) {
			tmp, err := menuService.SelectMenuNormalAll()
			if err == nil {
				menus = tmp
			}

		} else {
			tmp, err := menuService.SelectMenusByUserId(string(user.UserId))
			if err == nil {
				menus = tmp
			}
		}

		//获取配置数据
		sideTheme := configService.GetValueByKey("sys.index.sideTheme")
		skinName := configService.GetValueByKey("sys.index.skinName")
		response.BuildTpl(c, "index").WriteTpl(gin.H{
			"avatar":    avatar,
			"loginname": loginname,
			"username":  username,
			"menus":     menus,
			"sideTheme": sideTheme,
			"skinName":  skinName,
		})
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

//后台框架欢迎页面
func Main(c *gin.Context) {
	response.BuildTpl(c, "main").WriteTpl()
}

//下载文件
func Download(c *gin.Context) {
	fileName := c.Query("fileName")
	delete := c.Query("delete")

	if fileName == "" {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	// 创建路径
	curDir, err := os.Getwd()
	if err != nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "获取目录失败",
		})
		return
	}

	filepath := curDir + "/public/upload/" + fileName
	file, err := os.Open(filepath)

	defer file.Close()

	if err != nil {
		response.BuildTpl(c, model.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}

	b, _ := ioutil.ReadAll(file)
	c.Writer.Header().Add("Content-Disposition", "attachment")
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Write(b)

	if delete == "true" {
		os.Remove(filepath)
	}

}

//切换皮肤
func SwitchSkin(c *gin.Context) {
	response.BuildTpl(c, "skin").WriteTpl()
}

//注销
func Logout(c *gin.Context) {
	if userService.IsSignedIn(c) {
		userService.SignOut(c)
	}

	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}
