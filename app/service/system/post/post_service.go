package post

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
	postModel "yj-app/app/model/system/post"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*postModel.Entity, error) {
	entity := &postModel.Entity{PostId: id}
	_, err := entity.FindOne()
	return entity, err
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, err := (&postModel.Entity{PostId: id}).Delete()
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
	result, err := postModel.DeleteBatch(ida...)
	if err != nil {
		return 0
	}
	return result
}

//添加数据
func AddSave(req *postModel.AddReq, c *gin.Context) (int64, error) {
	var entity postModel.Entity
	entity.PostName = req.PostName
	entity.PostCode = req.PostCode
	entity.Status = req.Status
	entity.PostSort = req.PostSort
	entity.Remark = req.Remark
	entity.CreateTime = time.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(c)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	_, err := entity.Insert()
	return entity.PostId, err
}

//修改数据
func EditSave(req *postModel.EditReq, c *gin.Context) (int64, error) {
	entity := &postModel.Entity{PostId: req.PostId}
	ok, err := entity.FindOne()
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("数据不存在")
	}

	entity.PostName = req.PostName
	entity.PostCode = req.PostCode
	entity.Status = req.Status
	entity.Remark = req.Remark
	entity.PostSort = req.PostSort
	entity.UpdateTime = time.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(c)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	return entity.Update()
}

//根据条件分页查询角色数据
func SelectListAll(params *postModel.SelectPageReq) ([]postModel.EntityFlag, error) {
	return postModel.SelectListAll(params)
}

//根据条件分页查询角色数据
func SelectListByPage(params *postModel.SelectPageReq) ([]postModel.Entity, *page.Paging, error) {
	return postModel.SelectListByPage(params)
}

// 导出excel
func Export(param *postModel.SelectPageReq) (string, error) {
	head := []string{"岗位序号", "岗位名称", "岗位编码", "岗位排序", "状态"}
	col := []string{"post_id", "post_name", "post_code", "post_sort", "stat"}
	return postModel.SelectListExport(param, head, col)
}

//根据用户ID查询岗位
func SelectPostsByUserId(userId int64) ([]postModel.EntityFlag, error) {
	var paramsPost *postModel.SelectPageReq
	postAll, err := postModel.SelectListAll(paramsPost)

	if err != nil || postAll == nil {
		return nil, errors.New("未查询到岗位数据")
	}

	userPost, err := postModel.SelectPostsByUserId(userId)

	if err != nil || userPost == nil {
		return nil, errors.New("未查询到用户岗位数据")
	} else {
		for i := range postAll {
			for j := range userPost {
				if userPost[j].PostId == postAll[i].PostId {
					postAll[i].Flag = true
					break
				}
			}
		}
	}

	return postAll, nil
}

//检查角色名是否唯一
func CheckPostNameUniqueAll(postName string) string {
	post, err := postModel.CheckPostNameUniqueAll(postName)
	if err != nil {
		return "1"
	}
	if post != nil && post.PostId > 0 {
		return "1"
	}
	return "0"
}

//检查岗位名称是否唯一
func CheckPostNameUnique(postName string, postId int64) string {
	post, err := postModel.CheckPostNameUniqueAll(postName)
	if err != nil {
		return "1"
	}
	if post != nil && post.PostId > 0 && post.PostId != postId {
		return "1"
	}
	return "0"
}

//检查岗位编码是否唯一
func CheckPostCodeUniqueAll(postCode string) string {
	post, err := postModel.CheckPostCodeUniqueAll(postCode)
	if err != nil {
		return "1"
	}
	if post != nil && post.PostId > 0 {
		return "1"
	}
	return "0"
}

//检查岗位编码是否唯一
func CheckPostCodeUnique(postCode string, postId int64) string {
	post, err := postModel.CheckPostCodeUniqueAll(postCode)
	if err != nil {
		return "1"
	}
	if post != nil && post.PostId > 0 && post.PostId != postId {
		return "1"
	}
	return "0"
}
