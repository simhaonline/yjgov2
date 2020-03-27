package form

import (
	"github.com/gin-gonic/gin"
	"yj-app/app/yjgframe/response"
)

func Autocomplete(c *gin.Context) {

	response.BuildTpl(c, "demo/form/autocomplete").WriteTpl()
}

func Basic(c *gin.Context) {
	response.BuildTpl(c, "demo/form/basic").WriteTpl()
}

func Button(c *gin.Context) {
	response.BuildTpl(c, "demo/form/button").WriteTpl()
}

func Cards(c *gin.Context) {
	response.BuildTpl(c, "demo/form/cards").WriteTpl()
}

func Datetime(c *gin.Context) {
	response.BuildTpl(c, "demo/form/datetime").WriteTpl()
}

func Duallistbox(c *gin.Context) {
	response.BuildTpl(c, "demo/form/duallistbox").WriteTpl()
}

func Grid(c *gin.Context) {
	response.BuildTpl(c, "demo/form/grid").WriteTpl()
}

func Jasny(c *gin.Context) {
	response.BuildTpl(c, "demo/form/jasny").WriteTpl()
}

func Select(c *gin.Context) {
	response.BuildTpl(c, "demo/form/select").WriteTpl()
}

func Sortable(c *gin.Context) {
	response.BuildTpl(c, "demo/form/sortable").WriteTpl()
}

func Summernote(c *gin.Context) {
	response.BuildTpl(c, "demo/form/summernote").WriteTpl()
}

func Tabs_panels(c *gin.Context) {
	response.BuildTpl(c, "demo/form/tabs_panels").WriteTpl()
}

func Timeline(c *gin.Context) {
	response.BuildTpl(c, "demo/form/timeline").WriteTpl()
}

func Upload(c *gin.Context) {
	response.BuildTpl(c, "demo/form/upload").WriteTpl()
}

func Validate(c *gin.Context) {
	response.BuildTpl(c, "demo/form/validate").WriteTpl()
}

func Wizard(c *gin.Context) {
	response.BuildTpl(c, "demo/form/wizard").WriteTpl()
}
