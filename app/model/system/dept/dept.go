package dept

import (
	"errors"
	"fmt"
	"strings"
	"yj-app/app/yjgframe/db"
	"yj-app/app/yjgframe/utils/gconv"
)

// Fill with you ideas below.

// Entity is the golang structure for table sys_dept.
type EntityExtend struct {
	Entity     `xorm:"extends"`
	ParentName string `json:"parentName"`
}

//分页请求参数
type SelectPageReq struct {
	ParentId  int64  `form:"parentId"`      //父部门ID
	DeptName  string `form:"deptName"`      //部门名称
	Status    string `form:"status"`        //状态
	BeginTime string `form:"beginTime"`     //开始时间
	EndTime   string `form:"endTime"`       //结束时间
	PageNum   int    `form:"pageNum"`       //当前页码
	PageSize  int    `form:"pageSize"`      //每页数
	SortName  string `form:"orderByColumn"` //排序字段
	SortOrder string `form:"isAsc"`         //排序方式
}

//新增页面请求参数
type AddReq struct {
	ParentId int64  `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
	OrderNum int    `form:"orderNum" binding:"required"`
	Leader   string `form:"leader"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
	Status   string `form:"status"`
}

//修改页面请求参数
type EditReq struct {
	DeptId   int64  `form:"deptId" binding:"required"`
	ParentId int64  `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
	OrderNum int    `form:"orderNum" binding:"required"`
	Leader   string `form:"leader"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
	Status   string `form:"status"`
}

//检查菜单名称请求参数
type CheckDeptNameReq struct {
	DeptId   int64  `form:"deptId"  binding:"required"`
	ParentId int64  `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
}

//检查菜单名称请求参数
type CheckDeptNameALLReq struct {
	ParentId int64  `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
}

//根据部门ID查询信息
func SelectDeptById(id int64) (*EntityExtend, error) {
	db := db.Instance().Engine()

	if db == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	var result EntityExtend
	model := db.Table(TableName()).Alias("d")
	model.Select("d.dept_id, d.parent_id, d.ancestors, d.dept_name, d.order_num, d.leader, d.phone, d.email, d.status,(select dept_name from sys_dept where dept_id = d.parent_id) parent_name")
	model.Where("d.dept_id = ?", id)
	_, err := model.Get(&result)
	return &result, err
}

//根据ID查询所有子部门
func SelectChildrenDeptById(deptId int64) []*Entity {
	db := db.Instance().Engine()

	if db == nil {
		return nil
	}
	var rs []*Entity
	db.Table(TableName()).Where("find_in_set(?, ancestors)", deptId).Find(&rs)
	return rs
}

//删除部门管理信息
func DeleteDeptById(deptId int64) int64 {
	var entity Entity
	entity.DeptId = deptId
	entity.DelFlag = "2"
	rs, err := entity.Update()
	if err != nil {
		return 0
	}
	return rs
}

//修改子元素关系
func UpdateDeptChildren(deptId int64, newAncestors, oldAncestors string) {
	deptList := SelectChildrenDeptById(deptId)

	if deptList == nil || len(deptList) <= 0 {
		return
	}

	for _, tmp := range deptList {
		tmp.Ancestors = strings.ReplaceAll(tmp.Ancestors, oldAncestors, newAncestors)
	}

	ancestors := " case dept_id"
	idStr := ""

	for _, dept := range deptList {
		ancestors += " when " + gconv.String(dept.DeptId) + " then " + dept.Ancestors
		if idStr == "" {
			idStr = gconv.String(dept.DeptId)
		} else {
			idStr += "," + gconv.String(dept.DeptId)
		}
	}

	ancestors += " end "
	db := db.Instance().Engine()

	if db == nil {
		return
	}

	rs, err := db.Table(TableName()).Where("dept_id in(?)", deptId).Update(map[string]interface{}{"ancestors": ancestors})
	fmt.Printf("修改了%v行 错误信息：%v", rs, err.Error())
}

//查询部门管理数据
func SelectDeptList(parentId int64, deptName, status string) ([]Entity, error) {
	var result []Entity
	db := db.Instance().Engine()
	if db == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	model := db.Table(TableName()).Alias("d").Where("d.del_flag = '0'")
	if parentId > 0 {
		model.Where("d.parent_id = ?", parentId)
	}
	if deptName != "" {
		model.Where("d.dept_name like ?", "%"+deptName+"%")
	}
	if status != "" {
		model.Where("d.status = ?", status)
	}
	model.OrderBy("d.parent_id, d.order_num")

	err := model.Find(&result)

	return result, err
}

//根据角色ID查询部门
func SelectRoleDeptTree(roleId int64) ([]string, error) {
	db := db.Instance().Engine()
	if db == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	model := db.Table(TableName()).Alias("d").Join("LEFT", []string{"sys_role_dept", "rd"}, "d.dept_id = rd.dept_id")
	model.Where("d.del_flag = '0'").Where("rd.role_id = ?", roleId)
	model.OrderBy("d.parent_id, d.order_num ")
	model.Select("concat(d.dept_id, d.dept_name) as dept_name")

	var result []string
	var rs []Entity
	err := model.Find(&result)
	if err == nil && rs != nil && len(rs) > 0 {
		for _, record := range rs {
			if record.DeptName != "" {
				result = append(result, record.DeptName)
			}
		}
	}
	return result, nil
}

//查询部门是否存在用户
func CheckDeptExistUser(deptId int64) bool {
	db := db.Instance().Engine()
	if db == nil {
		return false
	}

	num, _ := db.Table(TableName()).Where("dept_id = ? and del_flag = '0'", deptId).Count()

	if num > 0 {
		return true
	} else {
		return false
	}
}

//查询部门人数
func SelectDeptCount(deptId, parentId int64) int64 {
	db := db.Instance().Engine()
	if db == nil {
		return 0
	}

	result := int64(0)
	whereStr := "del_flag = '0'"
	if deptId > 0 {
		whereStr = whereStr + " and dept_id=" + gconv.String(deptId)
	}
	if parentId > 0 {
		whereStr = whereStr + " and parent_id=" + gconv.String(parentId)
	}

	rs, err := db.Table(TableName()).Where(whereStr).Count()
	if err != nil {
		result = rs
	}
	return result
}

//校验部门名称是否唯一
func CheckDeptNameUniqueAll(deptName string, parentId int64) (*Entity, error) {
	var entity Entity
	entity.DeptName = deptName
	entity.ParentId = parentId
	ok, err := entity.FindOne()
	if ok {
		return &entity, err
	} else {
		return nil, err
	}
}
