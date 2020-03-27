package table

import (
	"time"
	"yj-app/app/yjgframe/db"
)

type Entity struct {
	TableId        int64     `json:"table_id" xorm:"not null pk autoincr comment('编号') BIGINT(20)"`
	TableName      string    `json:"table_name" xorm:"default '' comment('表名称') VARCHAR(200)"`
	TableComment   string    `json:"table_comment" xorm:"default '' comment('表描述') VARCHAR(500)"`
	ClassName      string    `json:"class_name" xorm:"default '' comment('实体类名称') VARCHAR(100)"`
	TplCategory    string    `json:"tpl_category" xorm:"default 'crud' comment('使用的模板（crud单表操作 tree树表操作）') VARCHAR(200)"`
	PackageName    string    `json:"package_name" xorm:"comment('生成包路径') VARCHAR(100)"`
	ModuleName     string    `json:"module_name" xorm:"comment('生成模块名') VARCHAR(30)"`
	BusinessName   string    `json:"business_name" xorm:"comment('生成业务名') VARCHAR(30)"`
	FunctionName   string    `json:"function_name" xorm:"comment('生成功能名') VARCHAR(50)"`
	FunctionAuthor string    `json:"function_author" xorm:"comment('生成功能作者') VARCHAR(50)"`
	Options        string    `json:"options" xorm:"comment('其它生成选项') VARCHAR(1000)"`
	CreateBy       string    `json:"create_by" xorm:"default '' comment('创建者') VARCHAR(64)"`
	CreateTime     time.Time `json:"create_time" xorm:"comment('创建时间') DATETIME"`
	UpdateBy       string    `json:"update_by" xorm:"default '' comment('更新者') VARCHAR(64)"`
	UpdateTime     time.Time `json:"update_time" xorm:"comment('更新时间') DATETIME"`
	Remark         string    `json:"remark" xorm:"comment('备注') VARCHAR(500)"`
}

//映射数据表
func TableName() string {
	return "gen_table"
}

// 插入数据
func (r *Entity) Insert() (int64, error) {
	return db.Instance().Engine().Table(TableName()).Insert(r)
}

// 更新数据
func (r *Entity) Update() (int64, error) {
	return db.Instance().Engine().Table(TableName()).ID(r.TableId).Update(r)
}

// 删除
func (r *Entity) Delete() (int64, error) {
	return db.Instance().Engine().Table(TableName()).ID(r.TableId).Delete(r)
}

//批量删除
func DeleteBatch(ids ...int64) (int64, error) {
	return db.Instance().Engine().Table(TableName()).In("table_id", ids).Delete(new(Entity))
}

// 根据结构体中已有的非空数据来获得单条数据
func (r *Entity) FindOne() (bool, error) {
	return db.Instance().Engine().Table(TableName()).Get(r)
}

// 根据条件查询
func Find(where, order string) ([]Entity, error) {
	var list []Entity
	err := db.Instance().Engine().Table(TableName()).Where(where).OrderBy(order).Find(&list)
	return list, err
}

//指定字段集合查询
func FindIn(column string, args ...interface{}) ([]Entity, error) {
	var list []Entity
	err := db.Instance().Engine().Table(TableName()).In(column, args).Find(&list)
	return list, err
}

//排除指定字段集合查询
func FindNotIn(column string, args ...interface{}) ([]Entity, error) {
	var list []Entity
	err := db.Instance().Engine().Table(TableName()).NotIn(column, args).Find(&list)
	return list, err
}
