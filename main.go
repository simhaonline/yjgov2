package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "yj-app/app/controller/router"
	"yj-app/app/service/middleware/sessions"
	"yj-app/app/service/middleware/sessions/memstore"
	"yj-app/app/yjgframe/cfg"
	_ "yj-app/app/yjgframe/cron"
	"yj-app/app/yjgframe/server"
)

var allDoneChan = make(chan struct{}, 1000)

// @title 云捷GO 自动生成API文档
// @version 1.0
// @description 生成文档请在调试模式下进行<a href="/tool/swagger?a=r">重新生成文档</a>

// @host localhost
// @BasePath /api
func main() {
	gin.SetMode("debug")
	config := cfg.Instance()

	if config == nil {
		fmt.Printf("参数错误")
		return
	}

	//后台服务状态
	adminStatus := config.Status.Admin
	//api服务状态
	apiStatus := config.Status.Api

	if adminStatus {
		store := memstore.NewStore([]byte("secret"))
		admin := server.New("admin", config.Admin.Address, gin.Logger(), sessions.Sessions("mysession", store))
		admin.Template("template").Static(config.Admin.ServerRoot)
		admin.Start()
	}

	if apiStatus {
		api := server.New("api", config.Api.Address, gin.Recovery(), gin.Logger())
		api.Start()
	}

	<-allDoneChan
}
