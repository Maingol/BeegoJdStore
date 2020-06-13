package controllers

import (
	"JDStore/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type GoodsController struct {
	beego.Controller
}

// 获取商品分类数据列表
func (this *GoodsController) GetGoodsCate() {
	// 不带分页参数返回的数据格式
	var resCateList models.ResCateList

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 149)
	if !hasRight {
		resCateList.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resCateList
		this.ServeJSON()
		return
	}

	// 获取数据
	typ, err := this.GetInt("type")
	errStr := fmt.Sprintf("%s", err)
	emptyErr := "strconv.Atoi: parsing \"\": invalid syntax"
	// 没有指定type参数时，将该参数赋值为0
	if errStr == emptyErr {
		typ = 0
	}
	// 只有不指定type或type为0，1，2，3时能通过检验
	if err != nil && errStr != emptyErr || typ != 0 && typ != 1 && typ != 2 && typ != 3 {
		resCateList.Meta = &models.ResMeta{"type参数错误", 400}
		this.Data["json"] = resCateList
		this.ServeJSON()
		return
	}

	pagenum, err := this.GetInt("pagenum")
	errStr = fmt.Sprintf("%s", err)
	if errStr == emptyErr {
		pagenum = 0
	}
	if err != nil && errStr != emptyErr || pagenum < 0 {
		resCateList.Meta = &models.ResMeta{"pagenum参数错误", 400}
		this.Data["json"] = resCateList
		this.ServeJSON()
		return
	}

	pagesize, err := this.GetInt("pagesize")
	errStr = fmt.Sprintf("%s", err)
	if errStr == emptyErr {
		pagesize = 0
	}
	if err != nil && errStr != emptyErr || pagesize < 0 {
		resCateList.Meta = &models.ResMeta{"pagesize参数错误", 400}
		this.Data["json"] = resCateList
		this.ServeJSON()
		return
	}

	if pagenum == 0 || pagesize == 0 {
		// 获取不带分页参数的商品分类列表
		resCateList.Data, err = models.GetCateList(typ)
		if err != nil {
			resCateList.Meta = &models.ResMeta{"获取不带分页参数的商品分类列表时出错", 400}
			this.Data["json"] = resCateList
			this.ServeJSON()
			return
		}
		resCateList.Meta = &models.ResMeta{"获取商品分类列表成功", 200}
		this.Data["json"] = resCateList
		this.ServeJSON()
		return
	}

	// 带分页参数返回的数据格式
	var resCatePage models.ResCatePage
	// 获取带分页参数的商品分类列表
	resCatePage.Data, err = models.GetCatePage(typ, pagenum, pagesize)
	if err != nil {
		resCatePage.Meta = &models.ResMeta{"获取带分页参数的商品分类列表时出错", 400}
		this.Data["json"] = resCatePage
		this.ServeJSON()
		return
	}
	resCatePage.Meta = &models.ResMeta{"获取商品分类列表成功", 200}
	this.Data["json"] = resCatePage
	this.ServeJSON()
}

// 添加商品分类
func (this *GoodsController) AddGoodsCate() {
	var resAddCate models.ResAddCate

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 122)
	if !hasRight {
		resAddCate.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 获取数据
	var addCateParams models.AddCateParams
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &addCateParams)
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"请求参数错误", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 校验请求参数
	valid := validation.Validation{}
	b, err := valid.Valid(&addCateParams)
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}
	if !b {
		msg := ""
		// 出错就退出循环，不必再继续遍历
		for _, err := range valid.Errors {
			//logs.Error(err.Field, err.Name, err.Key, err.Message)
			// err.Field是出错字段，err.Name是校验方法，err.Key是出错字段+校验方法，err.Message是错误信息
			msg = fmt.Sprintf("错误字段：%s，错误信息：%s", err.Field, err.Message)
			break
		}
		resAddCate.Meta = &models.ResMeta{msg, 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 数据校验通过，向数据库中插入数据
	category := models.SpCategory{
		CatPid:   addCateParams.Cat_pid,
		CatName:  addCateParams.Cat_name,
		CatLevel: addCateParams.Cat_level}
	o := orm.NewOrm()
	_, err = o.Insert(&category)
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"添加执行错误", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 添加分类成功，返回成功信息
	cate := models.GoodsCate{
		category.CatId, category.CatName,
		category.CatPid, category.CatLevel,
		category.CatDeleted == 1}
	resAddCate.Data = &cate
	resAddCate.Meta = &models.ResMeta{"添加分类成功", 201}
	this.Data["json"] = resAddCate
	this.ServeJSON()
}

// 修改分类名称
func (this *GoodsController) UpdateCateName() {
	var resAddCate models.ResAddCate
	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	o := orm.NewOrm()
	cate := models.SpCategory{CatId: id}
	err = o.Read(&cate, "CatId")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	var updateCateParams models.UpdateCateParams
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &updateCateParams)
	if err != nil {
		logs.Error("请求体中的参数解析错误：", err)
	}
	// 检验分类名称是否为空
	if updateCateParams.Cat_name == "" {
		resAddCate.Meta = &models.ResMeta{"分类名称不能为空", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 数据校验完成，修改分类名称
	cate.CatName = updateCateParams.Cat_name
	_, err = o.Update(&cate, "CatName")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"修改分类名称出错", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 修改完成，返回成功信息
	resCate := models.GoodsCate{
		cate.CatId, cate.CatName,
		cate.CatPid, cate.CatLevel,
		cate.CatDeleted == 1}
	resAddCate.Data = &resCate
	resAddCate.Meta = &models.ResMeta{"分类名称修改成功", 200}
	this.Data["json"] = resAddCate
	this.ServeJSON()
}

// 删除分类名称
func (this *GoodsController) DeleteCate() {
	var resAddCate models.ResAddCate

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 123)
	if !hasRight {
		resAddCate.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	o := orm.NewOrm()
	cate := models.SpCategory{CatId: id}
	err = o.Read(&cate, "CatId")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}
	// 数据校验通过，删除分类（实际上执行的是假删操作，修改删除状态）
	cate.CatDeleted = 1
	_, err = o.Update(&cate, "CatDeleted")
	if err != nil {
		resAddCate.Meta = &models.ResMeta{"修改执行出错", 400}
		this.Data["json"] = resAddCate
		this.ServeJSON()
		return
	}
	// 删除成功，返回成功信息
	resAddCate.Meta = &models.ResMeta{"删除分类成功", 200}
	this.Data["json"] = resAddCate
	this.ServeJSON()
}
