package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// 数据库中sp_category表的模型
type SpCategory struct {
	CatId      int `orm:"pk;auto"`
	CatName    string
	CatPid     int
	CatLevel   int
	CatDeleted int
	CatIcon    string
	CatSrc     string
}

// 定义商品分类中每个分类的结构
type GoodsCate struct {
	CatId      int    `json:"cat_id"`
	CatName    string `json:"cat_name"`
	CatPid     int    `json:"cat_pid"`
	CatLevel   int    `json:"cat_level"`
	CatDeleted bool   `json:"cat_deleted"`
}

// 定义商品分类返回结果中每个分类的结构
type ResGoodsCate struct {
	CatId      int             `json:"cat_id"`
	CatName    string          `json:"cat_name"`
	CatPid     int             `json:"cat_pid"`
	CatLevel   int             `json:"cat_level"`
	CatDeleted bool            `json:"cat_deleted"`
	Children   []*ResGoodsCate `json:"children"`
}

// 获取商品分类列表时，包含分页参数的返回结果中data的结构
type ResCatePageData struct {
	Total    int             `json:"total"`
	PageNum  int             `json:"pagenum"`
	PageSize int             `json:"pagesize"`
	Result   []*ResGoodsCate `json:"result"`
}

// 获取商品分类列表时，包含分页参数的返回结果
type ResCatePage struct {
	Data *ResCatePageData `json:"data"`
	Meta *ResMeta         `json:"meta"`
}

// 获取商品分类列表时，不包含分页参数的返回结果
type ResCateList struct {
	Data []*ResGoodsCate `json:"data"`
	Meta *ResMeta        `json:"meta"`
}

// post请求中请求参数的结构体 接口：categories
type AddCateParams struct {
	Cat_pid   int
	Cat_name  string `valid:"Required"`
	Cat_level int    `valid:"Range(0, 2)"`
}

// 添加分类接口返回的数据的结构
type ResAddCate struct {
	Data *GoodsCate `json:"data"`
	Meta *ResMeta   `json:"meta"`
}

// put请求中请求参数的结构体 接口：categories/:id
type UpdateCateParams struct {
	Cat_name string
}

// 请求参数CatPid的自定义校验规则
func (this *AddCateParams) Valid(v *validation.Validation) {
	var category SpCategory
	o := orm.NewOrm()
	category = SpCategory{CatId: this.Cat_pid}
	if this.Cat_pid != 0 {
		err := o.Read(&category, "CatId")
		if err != nil {
			v.SetError("CatPid", "分类父ID不存在")
		}
	}
	category = SpCategory{CatName: this.Cat_name}
	err := o.Read(&category, "CatName")
	if err == nil {
		logs.Info(category)
		v.SetError("CatName", "分类名称已存在")
	}
}

// 获取商品分类列表，不包含分页参数
func GetCateList(level int) ([]*ResGoodsCate, error) {
	return GetCate(0, level)
}

// 获取商品分类列表，包含分页参数
func GetCatePage(level int, pageNum int, pageSize int) (*ResCatePageData, error) {
	o := orm.NewOrm()
	// 过滤掉被假删的分类
	//total, err := o.QueryTable("sp_category").
	//	Filter("CatPid", 0).Exclude("CatDeleted", 1).Count()

	// 保留被假删的分类
	total, err := o.QueryTable("sp_category").
		Filter("CatPid", 0).Count()
	if err != nil {
		return nil, err
	}
	resCatePageData := ResCatePageData{Total: int(total), PageNum: pageNum, PageSize: pageSize}

	start := (pageNum - 1) * pageSize
	sqlStr := fmt.Sprintf(`
		SELECT
			cat_id, cat_name, cat_pid, cat_level, cat_deleted
		FROM
			sp_category
		WHERE
			cat_pid = 0
		LIMIT
			%d,%d`, start, pageSize)
	/* 存储从数据库中查询到的结果集 */
	var goodsCate []GoodsCate
	/* 从数据库中查询出所有pid为id的记录，例如id为0，查询结果就是所有一级权限 */
	_, err = o.Raw(sqlStr).QueryRows(&goodsCate)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resGoodsCats []*ResGoodsCate
	/* 遍历查询结果，结果中的每个字段都赋值给新的结构体，并且都添加一个children字段赋给这个新结构体 */
	for _, v := range goodsCate {
		children, err := GetCate(v.CatId, level)
		if err != nil {
			return nil, err
		}
		resGoodsCate := ResGoodsCate{
			CatId: v.CatId, CatName: v.CatName, CatPid: v.CatPid,
			CatLevel: v.CatLevel, CatDeleted: v.CatDeleted}
		resGoodsCate.Children = children
		resGoodsCats = append(resGoodsCats, &resGoodsCate)
	}
	resCatePageData.Result = resGoodsCats
	return &resCatePageData, nil
}

// 递归获取商品分类列表。第一个参数是父分类的id，第二个参数是递归截止到哪一级分类，第三个参数表示是否分页展示
func GetCate(id int, level int) ([]*ResGoodsCate, error) {
	query := map[int]string{1: " && cat_level = 0", 2: " && cat_level != 2", 3: "", 0: ""}

	sqlStr := fmt.Sprintf(`
		SELECT
			cat_id, cat_name, cat_pid, cat_level, cat_deleted
		FROM
			sp_category
		WHERE
			cat_pid = %d%s`, id, query[level])
	/* 存储从数据库中查询到的结果集 */
	var goodsCate []GoodsCate
	o := orm.NewOrm()
	/* 从数据库中查询出所有pid为id的记录，例如id为0，查询结果就是所有一级权限 */
	num, err := o.Raw(sqlStr).QueryRows(&goodsCate)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resGoodsCats []*ResGoodsCate
	/* 递归终止条件 */
	if num == 0 {
		/* 返回一个空切片 */
		return []*ResGoodsCate{}, nil
	}
	/* 遍历查询结果，结果中的每个字段都赋值给新的结构体，并且都添加一个children字段赋给这个新结构体 */
	for _, v := range goodsCate {
		children, err := GetCate(v.CatId, level)
		if err != nil {
			return nil, err
		}
		resGoodsCate := ResGoodsCate{
			CatId: v.CatId, CatName: v.CatName, CatPid: v.CatPid,
			CatLevel: v.CatLevel, CatDeleted: v.CatDeleted}
		resGoodsCate.Children = children
		resGoodsCats = append(resGoodsCats, &resGoodsCate)
	}
	return resGoodsCats, nil
}

//初始化模型
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(SpCategory))
}
