package table

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
	tableModel "yj-app/app/model/tool/table"
	tableColumnModel "yj-app/app/model/tool/table_column"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/cfg"
	"yj-app/app/yjgframe/db"
	"yj-app/app/yjgframe/utils/convert"
	"yj-app/app/yjgframe/utils/gconv"
	"yj-app/app/yjgframe/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*tableModel.EntityExtend, error) {
	entity, err := tableModel.SelectRecordById(id)
	if err != nil {
		return nil, err
	}
	//表附加属性
	SetTableFromOptions(entity)
	return entity, nil
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	rs, _ := (&tableModel.Entity{TableId: id}).Delete()
	if rs > 0 {
		return true
	}
	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, err := tableModel.DeleteBatch(idarr...)
	if err != nil {
		return 0
	}

	if result > 0 {
		db.Instance().Engine().SQL("delete from gen_table where table_id in (?)", idarr)
	}

	return result
}

//保存修改数据
func SaveEdit(req *tableModel.EditReq, c *gin.Context) (int64, error) {
	if req == nil {
		return 0, errors.New("参数错误")
	}

	table := tableModel.Entity{TableId: req.TableId}

	ok, err := table.FindOne()
	if err != nil || !ok {
		return 0, errors.New("数据不存在")
	}

	if req.TableName != "" {
		table.TableName = req.TableName
	}

	if req.TableComment != "" {
		table.TableComment = req.TableComment
	}

	if req.BusinessName != "" {
		table.BusinessName = req.BusinessName
	}

	if req.ClassName != "" {
		table.ClassName = req.ClassName
	}

	if req.FunctionAuthor != "" {
		table.FunctionAuthor = req.FunctionAuthor
	}

	if req.FunctionName != "" {
		table.FunctionName = req.FunctionName
	}

	if req.ModuleName != "" {
		table.ModuleName = req.ModuleName
	}

	if req.PackageName != "" {
		table.PackageName = req.PackageName
	}

	if req.Remark != "" {
		table.Remark = req.Remark
	}

	if req.TplCategory != "" {
		table.TplCategory = req.TplCategory
	}

	if req.Params != "" {
		table.Options = req.Params
	}

	table.UpdateTime = time.Now()

	user := userService.GetProfile(c)

	if user != nil {
		table.UpdateBy = user.LoginName
	}

	session := db.Instance().Engine().NewSession()

	tanErr := session.Begin()

	_, tanErr = session.Table(tableModel.TableName()).ID(table.TableId).Update(table)

	if tanErr != nil {
		session.Rollback()
		return 0, err
	}

	//保存列数据
	if req.Columns != "" {
		var columnList []tableColumnModel.Entity
		if err := json.Unmarshal([]byte(req.Columns), &columnList); err != nil {
			return 0, err
		} else {
			if err == nil && columnList != nil && len(columnList) > 0 {
				for _, column := range columnList {
					if column.ColumnId > 0 {
						tmp := new(tableColumnModel.Entity)
						tmp.ColumnId = column.ColumnId
						ok, _ := tmp.FindOne()
						if ok {
							tmp.ColumnComment = column.ColumnComment
							tmp.GoType = column.GoType
							tmp.HtmlType = column.HtmlType
							tmp.QueryType = column.QueryType
							tmp.GoField = column.GoField
							tmp.DictType = column.DictType
							tmp.IsInsert = column.IsInsert
							tmp.IsEdit = column.IsEdit
							tmp.IsList = column.IsList
							tmp.IsQuery = column.IsQuery

							_, err = session.Table(tableColumnModel.TableName()).ID(tmp.ColumnId).Update(tmp)

							if err != nil {
								session.Rollback()
								return 0, err
							}
						}
					}
				}
			}
		}
	}

	return 1, session.Commit()
}

//设置代码生成其他选项值
func SetTableFromOptions(entity *tableModel.EntityExtend) {
	if entity != nil && entity.Options != "" {
		p := make(map[string]interface{}, 2)
		if e := json.Unmarshal([]byte(entity.Options), &p); e == nil {
			treeCode := p["treeCode"].(string)
			treeParentCode := p["treeParentCode"].(string)
			treeName := p["treeName"].(string)
			entity.TreeCode = treeCode
			entity.TreeParentCode = treeParentCode
			entity.TreeName = treeName
		}
	}

}

//设置主键列信息
func SetPkColumn(table *tableModel.EntityExtend, columns []tableColumnModel.Entity) {
	for _, column := range columns {
		if column.IsPk == "1" {
			table.PkColumn = column
			break
		}
	}
	if &(table.PkColumn) == nil {
		table.PkColumn = columns[0]
	}
}

//根据条件分页查询数据
func SelectListByPage(param *tableModel.SelectPageReq) ([]tableModel.Entity, *page.Paging, error) {
	return tableModel.SelectListByPage(param)
}

//查询据库列表
func SelectDbTableList(param *tableModel.SelectPageReq) ([]tableModel.Entity, *page.Paging, error) {
	return tableModel.SelectDbTableList(param)
}

//查询据库列表
func SelectDbTableListByNames(tableNames []string) ([]tableModel.Entity, error) {
	return tableModel.SelectDbTableListByNames(tableNames)
}

//根据table_id查询表列数据
func SelectGenTableColumnListByTableId(tableId int64) ([]tableColumnModel.Entity, error) {
	return tableColumnModel.SelectGenTableColumnListByTableId(tableId)
}

//查询据库列表
func SelectTableByName(tableName string) (*tableModel.Entity, error) {
	return tableModel.SelectTableByName(tableName)
}

//查询表ID业务信息
func SelectGenTableById(tableId int64) (*tableModel.Entity, error) {
	return tableModel.SelectGenTableById(tableId)
}

//查询表名称业务信息
func SelectGenTableByName(tableName string) (*tableModel.Entity, error) {
	return tableModel.SelectGenTableByName(tableName)
}

//导入表结构
func ImportGenTable(tableList *[]tableModel.Entity, operName string) error {
	if tableList != nil && operName != "" {
		session := db.Instance().Engine().NewSession()
		err := session.Begin()
		defer session.Close()

		for _, table := range *tableList {
			tableName := table.TableName
			InitTable(&table, operName)
			_, err = session.Table(tableModel.TableName()).Insert(&table)
			if err != nil {
				return err
			}

			if err != nil || table.TableId <= 0 {
				session.Rollback()
				return errors.New("保存数据失败")
			}

			// 保存列信息
			genTableColumns, err := tableColumnModel.SelectDbTableColumnsByName(tableName)

			if err != nil || len(genTableColumns) <= 0 {
				session.Rollback()
				return errors.New("获取列数据失败")
			}

			for _, column := range genTableColumns {
				InitColumnField(&column, &table)
				_, err = session.Table(tableColumnModel.TableName()).Insert(&column)
				if err != nil {
					session.Rollback()
					return errors.New("保存列数据失败")
				}
			}
		}
		return session.Commit()
	} else {
		return errors.New("参数错误")
	}
}

//初始化表信息
func InitTable(table *tableModel.Entity, operName string) {
	table.ClassName = ConvertClassName(table.TableName)
	table.PackageName = cfg.Instance().Gen.PackageName
	table.ModuleName = cfg.Instance().Gen.ModuleName
	table.BusinessName = GetBusinessName(table.TableName)
	table.FunctionName = strings.ReplaceAll(table.TableComment, "表", "")
	table.FunctionAuthor = cfg.Instance().Gen.Author
	table.CreateBy = operName
	table.TplCategory = "crud"
	table.CreateTime = time.Now()
}

//初始化列属性字段
func InitColumnField(column *tableColumnModel.Entity, table *tableModel.Entity) {
	dataType := GetDbType(column.ColumnType)
	columnName := column.ColumnName
	column.TableId = table.TableId
	column.CreateBy = table.CreateBy
	//设置字段名
	column.GoField = ConvertToCamelCase(columnName)
	column.HtmlField = ConvertToCamelCase1(columnName)

	if tableColumnModel.IsStringObject(dataType) {
		//字段为字符串类型
		column.GoType = "string"
		if strings.EqualFold(dataType, "text") || strings.EqualFold(dataType, "tinytext") || strings.EqualFold(dataType, "mediumtext") || strings.EqualFold(dataType, "longtext") {
			column.HtmlType = "textarea"
		} else {
			columnLength := GetColumnLength(column.ColumnType)
			if columnLength >= 500 {
				column.HtmlType = "textarea"
			} else {
				column.HtmlType = "input"
			}
		}
	} else if tableColumnModel.IsTimeObject(dataType) {
		//字段为时间类型
		column.GoType = "Time"
		column.HtmlType = "datatime"
	} else if tableColumnModel.IsNumberObject(dataType) {
		//字段为数字类型
		column.HtmlType = "input"
		// 如果是浮点型
		tmp := column.ColumnType
		if tmp == "float" || tmp == "double" {
			column.GoType = "float64"
		} else {
			start := strings.Index(tmp, "(")
			end := strings.Index(tmp, ")")
			result := tmp[start+1 : end]
			arr := strings.Split(result, ",")
			if len(arr) == 2 && gconv.Int(arr[1]) > 0 {
				column.GoType = "float64"
			} else if len(arr) == 1 && gconv.Int(arr[0]) <= 10 {
				column.GoType = "int"
			} else {
				column.GoType = "int64"
			}
		}

	}
	//新增字段
	if columnName == "create_by" || columnName == "create_time" || columnName == "update_by" || columnName == "update_time" {
		column.IsRequired = "0"
		column.IsInsert = "0"
	} else {
		column.IsRequired = "0"
		column.IsInsert = "1"
		if strings.Index(columnName, "name") >= 0 || strings.Index(columnName, "status") >= 0 {
			column.IsRequired = "1"
		}
	}

	// 编辑字段
	if tableColumnModel.IsNotEdit(columnName) {
		if column.IsPk == "1" {
			column.IsEdit = "0"
		} else {
			column.IsEdit = "1"
		}
	} else {
		column.IsEdit = "0"
	}
	// 列表字段
	if tableColumnModel.IsNotList(columnName) {
		column.IsList = "1"
	} else {
		column.IsList = "0"
	}
	// 查询字段
	if tableColumnModel.IsNotQuery(columnName) {
		column.IsQuery = "1"
	} else {
		column.IsQuery = "0"
	}

	// 查询字段类型
	if CheckNameColumn(columnName) {
		column.QueryType = "LIKE"
	} else {
		column.QueryType = "EQ"
	}

	// 状态字段设置单选框
	if CheckStatusColumn(columnName) {
		column.HtmlType = "radio"
	} else if CheckTypeColumn(columnName) || CheckSexColumn(columnName) {
		// 类型&性别字段设置下拉框
		column.HtmlType = "select"
	}
}

//检查字段名后3位是否是sex
func CheckSexColumn(columnName string) bool {
	if len(columnName) >= 3 {
		end := len(columnName)
		start := end - 3

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "sex" {
			return true
		}
	}
	return false
}

//检查字段名后4位是否是type
func CheckTypeColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "type" {
			return true
		}
	}
	return false
}

//检查字段名后4位是否是name
func CheckNameColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		tmp := columnName[start:end]

		if tmp == "name" {
			return true
		}
	}
	return false
}

//检查字段名后6位是否是status
func CheckStatusColumn(columnName string) bool {
	if len(columnName) >= 6 {
		end := len(columnName)
		start := end - 6

		if start <= 0 {
			start = 0
		}
		tmp := columnName[start:end]

		if tmp == "status" {
			return true
		}
	}

	return false
}

//获取数据库类型字段
func GetDbType(columnType string) string {
	if strings.Index(columnType, "(") > 0 {
		return columnType[0:strings.Index(columnType, "(")]
	} else {
		return columnType
	}
}

//表名转换成类名
func ConvertClassName(tableName string) string {
	autoRemovePre := cfg.Instance().Gen.AutoRemovePre
	tablePrefix := cfg.Instance().Gen.TablePrefix
	if autoRemovePre && tablePrefix != "" {
		searchList := strings.Split(tablePrefix, ",")
		for _, str := range searchList {
			tableName = strings.ReplaceAll(tableName, str, "")
		}
	}
	return tableName
}

//获取业务名
func GetBusinessName(tableName string) string {
	lastIndex := strings.LastIndex(tableName, "_")
	nameLength := len(tableName)
	businessName := tableName[lastIndex+1 : nameLength]
	return businessName
}

//将下划线大写方式命名的字符串转换为驼峰式。如果转换前的下划线大写方式命名的字符串为空，则返回空字符串。 例如：HELLO_WORLD->HelloWorld
func ConvertToCamelCase(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// 不含下划线，仅将首字母大写
		return strings.ToUpper(name[0:1]) + name[1:len(name)]
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
	}
	return result
}

////将下划线大写方式命名的字符串转换为驼峰式,首字母小写。如果转换前的下划线大写方式命名的字符串为空，则返回空字符串。 例如：HELLO_WORLD->helloWorld
func ConvertToCamelCase1(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// 不含下划线，原值返回
		return name
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		if result == "" {
			result = strings.ToLower(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		} else {
			result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		}
	}
	return result
}

//获取字段长度
func GetColumnLength(columnType string) int {
	start := strings.Index(columnType, "(")
	end := strings.Index(columnType, ")")
	result := columnType[start+1 : end-1]
	return gconv.Int(result)
}

//获取Go类别下拉框
func GoTypeTpl() string {
	return `<script id="goTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].goType'>
    <option value="int64" {{if goType==="int64"}}selected{{/if}}>int64</option>
    <option value="int" {{if goType==="int"}}selected{{/if}}>int</option>
    <option value="string" {{if goType==="string"}}selected{{/if}}>string</option>
    <option value="Time" {{if goType==="Time"}}selected{{/if}}>Time</option>
    <option value="float64" {{if goType==="float64"}}selected{{/if}}>float64</option>
    <option value="byte" {{if goType==="byte"}}selected{{/if}}>byte</option>
</select>
</div>
</script>`
}

//获取查询方式下拉框
func QueryTypeTpl() string {
	return `<script id="queryTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].queryType'>
    <option value="EQ" {{if queryType==="EQ"}}selected{{/if}}>=</option>
    <option value="NE" {{if queryType==="NE"}}selected{{/if}}>!=</option>
    <option value="GT" {{if queryType==="GT"}}selected{{/if}}>></option>
    <option value="GTE" {{if queryType==="GTE"}}selected{{/if}}>>=</option>
    <option value="LT" {{if queryType==="LT"}}selected{{/if}}><</option>
    <option value="LTE" {{if queryType==="LTE"}}selected{{/if}}><=</option>
    <option value="LIKE" {{if queryType==="LIKE"}}selected{{/if}}>Like</option>
    <option value="BETWEEN" {{if queryType==="BETWEEN"}}selected{{/if}}>Between</option>
</select>
</div>
</script>`
}

// 获取显示类型下拉框
func HtmlTypeTpl() string {
	return `<script id="htmlTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].htmlType'>
    <option value="input" {{if htmlType==="input"}}selected{{/if}}>文本框</option>
    <option value="textarea" {{if htmlType==="textarea"}}selected{{/if}}>文本域</option>
    <option value="select" {{if htmlType==="select"}}selected{{/if}}>下拉框</option>
    <option value="radio" {{if htmlType==="radio"}}selected{{/if}}>单选框</option>
    <option value="checkbox" {{if htmlType==="checkbox"}}selected{{/if}}>复选框</option>
    <option value="datetime" {{if htmlType==="datetime"}}selected{{/if}}>日期控件</option>
</select>
</div>
</script>`
}

//读取模板
func LoadTemplate(templateName string, data interface{}) (string, error) {
	cur, err := os.Getwd()
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(cur + "/template/" + templateName)
	if err != nil {
		return "", err
	}
	templateStr := string(b)

	tmpl, err := template.New(templateName).Parse(templateStr) //建立一个模板，内容是"hello, {{.}}"
	if err != nil {
		return "", nil
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, data) //将string与模板合成，变量name的内容会替换掉{{.}}
	if err != nil {
		return "", nil
	}
	return buffer.String(), nil
}
