package report

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/yjgframe/response"
)

func Echarts(c *gin.Context) {
	response.BuildTpl(c, "demo/report/echarts").WriteTpl()
}

func Metrics(c *gin.Context) {
	response.BuildTpl(c, "demo/report/metrics").WriteTpl()
}

func Peity(c *gin.Context) {
	response.BuildTpl(c, "demo/report/peity").WriteTpl()
}

func Sparkline(c *gin.Context) {
	response.BuildTpl(c, "demo/report/sparkline").WriteTpl()
}
