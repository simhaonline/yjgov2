package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"yj-app/app/model"
	menuService "yj-app/app/service/system/menu"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/router"
)

// 鉴权中间件，只有登录成功之后才能通过
func Auth(c *gin.Context) {

	//判断是否登陆
	if userService.IsSignedIn(c) {
		//根据url判断是否有权限
		url := c.Request.URL.Path
		strEnd := url[len(url)-1 : len(url)]
		if strings.EqualFold(strEnd, "/") {
			url = strings.TrimRight(url, "/")
		}
		//获取权限标识
		permission := router.FindPermission(url)
		if len(permission) > 0 {
			//获取用户信息
			user := userService.GetProfile(c)
			//获取用户菜单列表
			menus, err := menuService.SelectMenuNormalByUser(user.UserId)
			if err != nil {

				c.Redirect(http.StatusFound, "/500")
				c.Abort()
				return
			}

			if menus == nil {
				c.Redirect(http.StatusFound, "/500")
				c.Abort()
				return
			}

			hasPermission := false

			for i := range *menus {
				if strings.EqualFold((*menus)[i].Perms, permission) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				ajaxString := c.Request.Header.Get("X-Requested-With")
				if strings.EqualFold(ajaxString, "XMLHttpRequest") {
					c.JSON(http.StatusOK, model.CommonRes{
						Code: 403,
						Msg:  "您没有操作权限",
					})
					c.Abort()
				} else {
					c.Redirect(http.StatusFound, "/403")
					c.Abort()
				}
			}
		}

		c.Next()
	} else {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
