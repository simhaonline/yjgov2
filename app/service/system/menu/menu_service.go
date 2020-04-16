package menu

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"yj-app/app/model"
	menuModel "yj-app/app/model/system/menu"
	userService "yj-app/app/service/system/user"
	cache "yj-app/app/yjgframe/cache"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/gconv"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*menuModel.EntityExtend, error) {
	return menuModel.SelectRecordById(id)
}

//根据条件查询数据
func SelectListAll(params *menuModel.SelectPageReq) ([]menuModel.Entity, error) {
	return menuModel.SelectListAll(params)
}

//根据条件分页查询数据
func SelectListPage(params *menuModel.SelectPageReq) (*[]menuModel.Entity, *page.Paging, error) {
	return menuModel.SelectListPage(params)
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, err := (&menuModel.Entity{MenuId: id}).Delete()
	if err == nil {
		if rs > 0 {
			return true
		}
	}
	return false
}

//添加数据
func AddSave(req *menuModel.AddReq, c *gin.Context) (int64, error) {

	var entity menuModel.Entity
	entity.MenuName = req.MenuName
	entity.Visible = req.Visible
	entity.ParentId = req.ParentId
	entity.Remark = ""
	entity.MenuType = req.MenuType
	entity.Url = req.Url
	entity.Perms = req.Perms
	entity.Target = req.Target
	entity.Icon = req.Icon
	entity.OrderNum = req.OrderNum
	entity.CreateTime = time.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(c)

	if user == nil {
		entity.CreateBy = user.LoginName
	}

	_, err := entity.Insert()
	return entity.MenuId, err
}

//修改数据
func EditSave(req *menuModel.EditReq, c *gin.Context) (int64, error) {
	entity := &menuModel.Entity{MenuId: req.MenuId}
	ok, err := entity.FindOne()

	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("角色不存在")
	}

	entity.MenuName = req.MenuName
	entity.Visible = req.Visible
	entity.ParentId = req.ParentId
	entity.Remark = ""
	entity.MenuType = req.MenuType
	entity.Url = req.Url
	entity.Perms = req.Perms
	entity.Target = req.Target
	entity.Icon = req.Icon
	entity.OrderNum = req.OrderNum
	entity.UpdateTime = time.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(c)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	return entity.Update()
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, err := menuModel.DeleteBatch(idarr...)
	if err != nil {
		return 0
	}
	return result
}

//加载所有菜单列表树
func MenuTreeData(userId int64) (*[]model.Ztree, error) {
	var result *[]model.Ztree
	menuList, err := SelectMenuNormalByUser(userId)
	if err != nil {
		return nil, err
	}
	result, err = InitZtree(menuList, nil, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//获取用户的菜单数据
func SelectMenuNormalByUser(userId int64) (*[]menuModel.EntityExtend, error) {
	if userService.IsAdmin(userId) {
		return SelectMenuNormalAll()
	} else {
		return SelectMenusByUserId(gconv.String(userId))
	}
}

//获取管理员菜单数据
func SelectMenuNormalAll() (*[]menuModel.EntityExtend, error) {

	//从缓存读取
	c := cache.Instance()
	tmp, f := c.Get(model.MENU_CACHE)

	if f && tmp != nil {
		rs, ok := tmp.([]menuModel.EntityExtend)
		if ok {
			return &rs, nil
		}
	}

	//从数据库中读取
	var result []menuModel.EntityExtend
	result, err := menuModel.SelectMenuNormalAll()

	if err != nil {
		return nil, err
	}

	for i := range result {
		chilrens := getMenuChildPerms(result, result[i].MenuId)

		for j := range chilrens {
			chilrens2 := getMenuChildPerms(result, chilrens[j].MenuId)
			chilrens[j].Children = chilrens2

			if chilrens[j].Target == "" {
				chilrens[j].Target = "menuItem"
			}
			if chilrens[j].Url == "" {
				chilrens[j].Url = "#"
			}
		}

		if chilrens != nil {
			result[i].Children = chilrens

			if result[i].ParentId != 0 {
				if result[i].Target == "" {
					result[i].Target = "menuItem"
				}

				if result[i].Url == "" {
					result[i].Url = "#"
				}
			}

		}
	}

	//存入缓存
	cache.Instance().Set(model.MENU_CACHE, result, time.Hour)
	return &result, nil
}

//根据用户ID读取菜单数据
func SelectMenusByUserId(userId string) (*[]menuModel.EntityExtend, error) {
	//从缓存读取
	tmp, have := cache.Instance().Get(model.MENU_CACHE + userId)

	if have && tmp != nil {
		rs, ok := tmp.([]menuModel.EntityExtend)
		if ok {
			return &rs, nil
		}
	}

	//从数据库中读取
	var result []menuModel.EntityExtend
	result, err := menuModel.SelectMenusByUserId(userId)

	if err != nil {
		return nil, err
	}

	for i := range result {
		chilrens := getMenuChildPerms(result, result[i].MenuId)

		for j := range chilrens {
			chilrens2 := getMenuChildPerms(result, chilrens[j].MenuId)
			chilrens[j].Children = chilrens2
			if chilrens[j].Target == "" {
				chilrens[j].Target = "menuItem"
			}
			if chilrens[j].Url == "" {
				chilrens[j].Url = "#"
			}
		}

		if chilrens != nil {
			result[i].Children = chilrens
			if result[i].ParentId != 0 {
				if result[i].Target == "" {
					result[i].Target = "menuItem"
				}

				if result[i].Url == "" {
					result[i].Url = "#"
				}
			} else {
				if result[i].Url == "" || result[i].Url == "#" {
					result[i].Target = ""
				}
				if result[i].Url == "" {
					result[i].Url = "#"
				}
			}
		}
	}

	//存入缓存
	cache.Instance().Set(model.MENU_CACHE+userId, result, time.Hour)
	return &result, nil
}

//根据父id获取子菜单
func getMenuChildPerms(menus []menuModel.EntityExtend, parentId int64) []menuModel.EntityExtend {
	if menus == nil {
		return nil
	}

	var result []menuModel.EntityExtend
	//得到一级菜单
	for i := range menus {
		if menus[i].ParentId == parentId && (menus[i].MenuType == "M" || menus[i].MenuType == "C") {
			if menus[i].Target == "" {
				menus[i].Target = "menuItem"
			}

			if menus[i].Url == "" {
				menus[i].Url = "#"
			}

			result = append(result, menus[i])
		}
	}

	return result
}

//检查菜单名是否唯一
func CheckMenuNameUniqueAll(menuName string, parentId int64) string {
	entity, err := menuModel.CheckMenuNameUniqueAll(menuName, parentId)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.MenuId > 0 {
		return "1"
	}
	return "0"
}

//检查菜单名是否唯一
func CheckMenuNameUnique(menuName string, menuId, parentId int64) string {
	entity, err := menuModel.CheckMenuNameUniqueAll(menuName, parentId)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.MenuId > 0 && entity.MenuId != menuId {
		return "1"
	}
	return "0"
}

//检查权限键是否唯一
func CheckPermsUniqueAll(perms string) string {
	entity, err := menuModel.CheckPermsUniqueAll(perms)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.MenuId > 0 {
		return "1"
	}
	return "0"
}

//检查权限键是否唯一
func CheckPermsUnique(perms string, menuId int64) string {
	entity, err := menuModel.CheckPermsUniqueAll(perms)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.MenuId > 0 && entity.MenuId != menuId {
		return "1"
	}
	return "0"
}

//根据角色ID查询菜单
func RoleMenuTreeData(roleId, userId int64) (*[]model.Ztree, error) {
	var result *[]model.Ztree
	menuList, err := SelectMenuNormalByUser(userId)
	if err != nil {
		return nil, err
	}

	if roleId > 0 {
		roleMenuList, err := menuModel.SelectMenuTree(roleId)
		if err != nil || roleMenuList == nil {
			result, err = InitZtree(menuList, nil, true)
		} else {
			result, err = InitZtree(menuList, &roleMenuList, true)
		}
	} else {
		result, err = InitZtree(menuList, nil, true)
	}

	return result, nil
}

//对象转菜单树
func InitZtree(menuList *[]menuModel.EntityExtend, roleMenuList *[]string, permsFlag bool) (*[]model.Ztree, error) {
	var result []model.Ztree
	isCheck := false
	if roleMenuList != nil && len(*roleMenuList) > 0 {
		isCheck = true
	}

	for _, obj := range *menuList {
		var ztree model.Ztree
		ztree.Title = obj.MenuName
		ztree.Id = obj.MenuId
		ztree.Name = transMenuName(obj.MenuName, permsFlag)
		ztree.Pid = obj.ParentId
		if isCheck {
			tmp := gconv.String(obj.MenuId) + obj.Perms
			tmpcheck := false
			for j := range *roleMenuList {
				if strings.Compare((*roleMenuList)[j], tmp) == 0 {
					tmpcheck = true
					break
				}
			}
			ztree.Checked = tmpcheck
		}
		result = append(result, ztree)
	}

	return &result, nil
}

func transMenuName(menuName string, permsFlag bool) string {
	if permsFlag {
		return "<font color=\"#888\">&nbsp;&nbsp;&nbsp;" + menuName + "</font>"
	} else {
		return menuName
	}
}
