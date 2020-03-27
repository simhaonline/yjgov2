package modal

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/yjgframe/response"
)

func Dialog(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/dialog").WriteTpl()
}

func Form(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/form").WriteTpl()
}

func Layer(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/layer").WriteTpl()
}

func Table(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/table").WriteTpl()
}

func Check(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/table/check").WriteTpl()
}

func Parent(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/table/parent").WriteTpl()
}

func Radio(c *gin.Context) {
	response.BuildTpl(c, "demo/modal/table/radio").WriteTpl()
}
