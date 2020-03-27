package table_column

import (
	tableColumn "yj-app/app/model/tool/table_column"
	"yj-app/app/yjgframe/utils/convert"
)

//新增业务字段
func Insert(entity *tableColumn.Entity) (int64, error) {
	_, err := entity.Insert()
	if err != nil {
		return 0, err
	}
	return entity.ColumnId, err
}

//修改业务字段
func Update(entity *tableColumn.Entity) (int64, error) {
	return entity.Update()
}

//根据主键查询数据
func SelectRecordById(id int64) (*tableColumn.Entity, error) {
	entity := &tableColumn.Entity{ColumnId: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, err := (&tableColumn.Entity{ColumnId: id}).Delete()
	if err == nil && rs > 0 {
		return true
	}
	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, err := tableColumn.DeleteBatch(idarr...)
	if err != nil {
		return 0
	}
	return result
}

//查询业务字段列表
func SelectGenTableColumnListByTableId(tableId int64) (*[]tableColumn.Entity, error) {
	return tableColumn.SelectGenTableColumnListByTableId(tableId)
}

//根据表名称查询列信息
func SelectDbTableColumnsByName(tableName string) (*[]tableColumn.Entity, error) {
	return tableColumn.SelectDbTableColumnsByName(tableName)
}
