package dept

import (
	"time"
	"yj-app/app/yjgframe/db"
)

type Entity struct {
	DeptId     int64     `json:"dept_id" xorm:"not null pk autoincr comment('部门id') BIGINT(20)"`
	ParentId   int64     `json:"parent_id" xorm:"default 0 comment('父部门id') BIGINT(20)"`
	Ancestors  string    `json:"ancestors" xorm:"default '' comment('祖级列表') VARCHAR(50)"`
	DeptName   string    `json:"dept_name" xorm:"default '' comment('部门名称') VARCHAR(30)"`
	OrderNum   int       `json:"order_num" xorm:"default 0 comment('显示顺序') INT(4)"`
	Leader     string    `json:"leader" xorm:"comment('负责人') VARCHAR(20)"`
	Phone      string    `json:"phone" xorm:"comment('联系电话') VARCHAR(11)"`
	Email      string    `json:"email" xorm:"comment('邮箱') VARCHAR(50)"`
	Status     string    `json:"status" xorm:"default '0' comment('部门状态（0正常 1停用）') CHAR(1)"`
	DelFlag    string    `json:"del_flag" xorm:"default '0' comment('删除标志（0代表存在 2代表删除）') CHAR(1)"`
	CreateBy   string    `json:"create_by" xorm:"default '' comment('创建者') VARCHAR(64)"`
	CreateTime time.Time `json:"create_time" xorm:"comment('创建时间') DATETIME"`
	UpdateBy   string    `json:"update_by" xorm:"default '' comment('更新者') VARCHAR(64)"`
	UpdateTime time.Time `json:"update_time" xorm:"comment('更新时间') DATETIME"`
}

//映射数据表
func TableName() string {
	return "sys_dept"
}

// 插入数据
func (r *Entity) Insert() (int64, error) {
	return db.Instance().Engine().Table(TableName()).Insert(r)
}

// 更新数据
func (r *Entity) Update() (int64, error) {
	return db.Instance().Engine().Table(TableName()).ID(r.DeptId).Update(r)
}

// 删除
func (r *Entity) Delete() (int64, error) {
	return db.Instance().Engine().Table(TableName()).ID(r.DeptId).Delete(r)
}

//批量删除
func DeleteBatch(ids ...int64) (int64, error) {
	return db.Instance().Engine().Table(TableName()).In("dept_id", ids).Delete(new(Entity))
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
