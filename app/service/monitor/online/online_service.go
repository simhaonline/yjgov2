// ==========================================================================
// 云捷GO自动生成业务逻辑层相关代码，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2020-02-17 14:03:51
// 生成路径: app/service/module/online/online_service.go
// 生成人：yunjie
// ==========================================================================
package online

import (
	"strings"
	onlineModel "yj-app/app/model/monitor/online"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id string) (*onlineModel.Entity, error) {
	entity := &onlineModel.Entity{Sessionid: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id string) bool {
	entity := &onlineModel.Entity{Sessionid: id}
	result, err := entity.Delete()
	if err == nil && result > 0 {
		return true
	}

	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := strings.Split(ids, ",")
	result, _ := onlineModel.DeleteBatch(idarr...)
	return result
}

//批量删除数据
func DeleteRecordNotInIds(ids []string) int64 {
	result, _ := onlineModel.DeleteNotIn(ids...)
	return result
}

//添加数据
func AddSave(entity onlineModel.Entity) (int64, error) {
	return entity.Insert()
}

//根据条件查询数据
func SelectListAll(params *onlineModel.SelectPageReq) ([]onlineModel.Entity, error) {
	return onlineModel.SelectListAll(params)
}

//根据条件分页查询数据
func SelectListByPage(params *onlineModel.SelectPageReq) ([]onlineModel.Entity, *page.Paging, error) {
	return onlineModel.SelectListByPage(params)
}
