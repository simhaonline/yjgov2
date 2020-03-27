package dict_type

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
	"yj-app/app/model"
	dictTypeModel "yj-app/app/model/system/dict_type"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*dictTypeModel.Entity, error) {
	entity := &dictTypeModel.Entity{DictId: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, err := (&dictTypeModel.Entity{DictId: id}).Delete()
	if err == nil {
		if rs > 0 {
			return true
		}
	}
	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	ida := convert.ToInt64Array(ids, ",")
	result, err := dictTypeModel.DeleteBatch(ida...)
	if err != nil {
		return 0
	}
	return result
}

//添加数据
func AddSave(req *dictTypeModel.AddReq, c *gin.Context) (int64, error) {
	var entity dictTypeModel.Entity
	entity.Status = req.Status
	entity.DictType = req.DictType
	entity.DictName = req.DictName
	entity.Remark = req.Remark
	entity.CreateTime = time.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(c)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	_, err := entity.Insert()

	return entity.DictId, err
}

//修改数据
func EditSave(req *dictTypeModel.EditReq, c *gin.Context) (int64, error) {
	entity := &dictTypeModel.Entity{DictId: req.DictId}
	ok, err := entity.FindOne()

	if err != nil || !ok {
		return 0, err
	}

	if entity == nil {
		return 0, errors.New("数据不存在")
	}
	entity.Status = req.Status
	entity.DictType = req.DictType
	entity.DictName = req.DictName
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
func SelectListAll(params *dictTypeModel.SelectPageReq) ([]dictTypeModel.Entity, error) {
	return dictTypeModel.SelectListAll(params)
}

//根据条件分页查询角色数据
func SelectListByPage(params *dictTypeModel.SelectPageReq) ([]dictTypeModel.Entity, *page.Paging, error) {
	return dictTypeModel.SelectListByPage(params)
}

//根据字典类型查询信息
func SelectDictTypeByType(dictType string) *dictTypeModel.Entity {
	entity := &dictTypeModel.Entity{DictType: dictType}
	ok, err := entity.FindOne()
	if err != nil || !ok {
		return nil
	}
	return entity
}

// 导出excel
func Export(param *dictTypeModel.SelectPageReq) (string, error) {
	head := []string{"字典主键", "字典名称", "字典类型", "状态", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	col := []string{"dict_id", "dict_name", "dict_type", "status", "create_by", "create_time", "update_by", "update_time", "remark"}
	return dictTypeModel.SelectListExport(param, head, col)
}

//检查字典类型是否唯一
func CheckDictTypeUniqueAll(configKey string) string {
	entity, err := dictTypeModel.CheckDictTypeUniqueAll(configKey)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.DictId > 0 {
		return "1"
	}
	return "0"
}

//检查字典类型是否唯一
func CheckDictTypeUnique(configKey string, dictId int64) string {
	entity, err := dictTypeModel.CheckDictTypeUniqueAll(configKey)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.DictId > 0 && entity.DictId != dictId {
		return "1"
	}
	return "0"
}

//查询字典类型树
func SelectDictTree(params *dictTypeModel.SelectPageReq) *[]model.Ztree {
	var result []model.Ztree
	dictList, err := dictTypeModel.SelectListAll(params)
	if err == nil && dictList != nil {
		for _, item := range dictList {
			var tmp model.Ztree
			tmp.Id = item.DictId
			tmp.Name = transDictName(item)
			tmp.Title = item.DictType
			result = append(result, tmp)
		}
	}
	return &result
}

func transDictName(entity dictTypeModel.Entity) string {
	return `(` + entity.DictName + `)&nbsp;&nbsp;&nbsp;` + entity.DictType
}
