package table

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yj-app/app/model"
	"yj-app/app/yjgframe/response"
)

func Button(c *gin.Context) {
	response.BuildTpl(c, "demo/table/button").WriteTpl()
}

func Child(c *gin.Context) {
	response.BuildTpl(c, "demo/table/child").WriteTpl()
}

func Curd(c *gin.Context) {
	response.BuildTpl(c, "demo/table/curd").WriteTpl()
}

func Detail(c *gin.Context) {
	response.BuildTpl(c, "demo/table/detail").WriteTpl()
}

func Editable(c *gin.Context) {
	response.BuildTpl(c, "demo/table/editable").WriteTpl()
}

func Event(c *gin.Context) {
	response.BuildTpl(c, "demo/table/event").WriteTpl()
}

func Export(c *gin.Context) {
	response.BuildTpl(c, "demo/table/export").WriteTpl()
}

func FixedColumns(c *gin.Context) {
	response.BuildTpl(c, "demo/table/fixedColumns").WriteTpl()
}

func Footer(c *gin.Context) {
	response.BuildTpl(c, "demo/table/footer").WriteTpl()
}

func GroupHeader(c *gin.Context) {
	response.BuildTpl(c, "demo/table/groupHeader").WriteTpl()
}

func Image(c *gin.Context) {
	response.BuildTpl(c, "demo/table/image").WriteTpl()
}

func Multi(c *gin.Context) {
	response.BuildTpl(c, "demo/table/multi").WriteTpl()
}

func Other(c *gin.Context) {
	response.BuildTpl(c, "demo/table/other").WriteTpl()
}

func PageGo(c *gin.Context) {
	response.BuildTpl(c, "demo/table/pageGo").WriteTpl()
}

func Params(c *gin.Context) {
	response.BuildTpl(c, "demo/table/params").WriteTpl()
}

func Remember(c *gin.Context) {
	response.BuildTpl(c, "demo/table/remember").WriteTpl()
}

func Recorder(c *gin.Context) {
	response.BuildTpl(c, "demo/table/recorder").WriteTpl()
}

func Search(c *gin.Context) {
	response.BuildTpl(c, "demo/table/search").WriteTpl()
}

type us struct {
	UserId      int64   `json:"userId"`
	UserCode    string  `json:"userCode"`
	UserName    string  `json:"userName"`
	UserSex     string  `json:"userName"`
	UserPhone   string  `json:"userPhone"`
	UserEmail   string  `json:"userEmail"`
	UserBalance float64 `json:"userBalance"`
	Status      string  `json:"status"`
	CreateTime  string  `json:"createTime"`
}

func List(c *gin.Context) {
	var rows = make([]us, 0)
	for i := 1; i <= 10; i++ {
		var tmp us
		tmp.UserId = int64(i)
		tmp.UserName = "测试" + string(i)
		tmp.Status = "0"
		tmp.CreateTime = "2020-01-12 02:02:02"
		tmp.UserBalance = 100
		tmp.UserCode = "100000" + string(i)
		tmp.UserSex = "0"
		tmp.UserPhone = "15888888888"
		tmp.UserEmail = "111@qq.com"
		rows = append(rows, tmp)
	}
	c.JSON(http.StatusOK, model.TableDataInfo{
		Code:  0,
		Msg:   "操作成功",
		Total: len(rows),
		Rows:  rows,
	})
}
