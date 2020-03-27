package icon

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/yjgframe/response"
)

func Fontawesome(c *gin.Context) {
	response.BuildTpl(c, "demo/icon/fontawesome").WriteTpl()
}

func Glyphicons(c *gin.Context) {
	response.BuildTpl(c, "demo/icon/glyphicons").WriteTpl()
}
