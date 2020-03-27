package dict_data

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
	dictDataModel "yj-app/app/model/system/dict_data"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*dictDataModel.Entity, error) {
	entity := &dictDataModel.Entity{DictCode: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, _ := (&dictDataModel.Entity{DictCode: id}).Delete()
	if rs > 0 {
		return true
	}
	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	ida := convert.ToInt64Array(ids, ",")
	result, err := dictDataModel.DeleteBatch(ida...)
	if err != nil {
		return 0
	}
	return result
}

//添加数据
func AddSave(req *dictDataModel.AddReq, c *gin.Context) (int64, error) {
	var entity dictDataModel.Entity
	entity.DictType = req.DictType
	entity.Status = req.Status
	entity.DictLabel = req.DictLabel
	entity.CssClass = req.CssClass
	entity.DictSort = req.DictSort
	entity.DictValue = req.DictValue
	entity.IsDefault = req.IsDefault
	entity.ListClass = req.ListClass
	entity.Remark = req.Remark
	entity.CreateTime = time.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(c)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	_, err := entity.Insert()

	return entity.DictCode, err
}

//修改数据
func EditSave(req *dictDataModel.EditReq, c *gin.Context) (int64, error) {
	entity := &dictDataModel.Entity{DictCode: req.DictCode}
	ok, err := entity.FindOne()

	if err != nil || !ok {
		return 0, err
	}

	if entity == nil {
		return 0, errors.New("数据不存在")
	}

	entity.DictType = req.DictType
	entity.Status = req.Status
	entity.DictLabel = req.DictLabel
	entity.CssClass = req.CssClass
	entity.DictSort = req.DictSort
	entity.DictValue = req.DictValue
	entity.IsDefault = req.IsDefault
	entity.ListClass = req.ListClass
	entity.Remark = req.Remark
	entity.UpdateTime = time.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(c)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	return entity.Update()
}

//根据条件分页查询角色数据
func SelectListAll(params *dictDataModel.SelectPageReq) ([]dictDataModel.Entity, error) {
	return dictDataModel.SelectListAll(params)
}

//根据条件分页查询角色数据
func SelectListByPage(params *dictDataModel.SelectPageReq) (*[]dictDataModel.Entity, *page.Paging, error) {
	return dictDataModel.SelectListByPage(params)
}

// 导出excel
func Export(param *dictDataModel.SelectPageReq) (string, error) {
	head := []string{"字典编码", "字典排序", "字典标签", "字典键值", "字典类型", "样式属性", "表格回显样式", "是否默认", "状态", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	col := []string{"dict_code", "dict_sort", "dict_label", "dict_value", "dict_type", "css_class", "list_class", "is_default", "status", "create_by", "create_time", "update_by", "update_time", "remark"}
	return dictDataModel.SelectListExport(param, head, col)
}
