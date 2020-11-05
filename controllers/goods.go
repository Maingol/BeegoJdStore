package controllers

import (
	"JDStore/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/rs/xid"
	"path"
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

// 获取参数列表
func (this *GoodsController) GetAttrList() {
	var resAttrList models.ResAttrList

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 142)
	if !hasRight {
		logs.Error("权限不足")
		resAttrList.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		logs.Error("分类id格式错误")
		resAttrList.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	// 获取分类id
	/*o := orm.NewOrm()
	cate := models.SpCategory{CatId: id}
	err = o.Read(&cate, "CatId")
	if err != nil {
		resAttrList.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}*/
	// 检验分类id是否存在
	if !models.CateIdExist(id) {
		logs.Error("分类id不存在")
		resAttrList.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}

	// 获取属性类型
	sel := this.GetString("sel")
	if sel == "" {
		logs.Error("属性类型不能为空")
		resAttrList.Meta = &models.ResMeta{"属性类型不能为空", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	} else if sel != "only" && sel != "many" {
		logs.Error("属性类型必须是'only'或者'many'")
		resAttrList.Meta = &models.ResMeta{"属性类型必须是'only'或者'many'", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}

	// 数据校验完成，执行查询
	var attrs []*models.SpAttribute
	o := orm.NewOrm()
	_, err = o.
		QueryTable("sp_attribute").
		Filter("cat_id", id).
		Filter("delete_time__isnull", true).
		Filter("attr_sel", sel).
		All(&attrs)
	if err != nil {
		logs.Error("查询参数列表出错")
		resAttrList.Meta = &models.ResMeta{"查询参数列表出错", 400}
		this.Data["json"] = resAttrList
		this.ServeJSON()
		return
	}

	// 查询执行完成，返回数据
	resAttrList.Data = attrs
	resAttrList.Meta = &models.ResMeta{"获取参数列表成功", 200}
	this.Data["json"] = resAttrList
	this.ServeJSON()
}

// 添加动态参数或者静态属性 接口：categories/:id/attributes 请求方式：post
func (this *GoodsController) AddAttr() {
	var resAttr models.ResAttr

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 156)
	if !hasRight {
		logs.Error("权限不足")
		resAttr.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		logs.Error("分类id格式错误")
		resAttr.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	if !models.CateIdExist(id) {
		logs.Error("分类id不存在")
		resAttr.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	var addAttrParams models.AddAttrParams
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &addAttrParams)
	if err != nil {
		logs.Error("参数解析错误")
	}
	if addAttrParams.Attr_name == "" {
		logs.Error("参数名称不能为空")
		resAttr.Meta = &models.ResMeta{"参数名称不能为空", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	if addAttrParams.Attr_sel == "" {
		logs.Error("属性类型不能为空")
		resAttr.Meta = &models.ResMeta{"属性类型不能为空", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	} else if addAttrParams.Attr_sel != "only" && addAttrParams.Attr_sel != "many" {
		logs.Error("属性类型必须是'only'或者'many'")
		resAttr.Meta = &models.ResMeta{"属性类型必须是'only'或者'many'", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 验证某一分类下的分类名称是否存在
	if models.AttrExist(id, addAttrParams.Attr_name, addAttrParams.Attr_sel) {
		logs.Error("分类参数不可重复添加")
		resAttr.Meta = &models.ResMeta{"分类参数不可重复添加", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 数据验证完毕，执行数据库添加操作
	attr, err := models.AddAttr(id, addAttrParams.Attr_name, addAttrParams.Attr_sel, addAttrParams.Attr_vals)
	if err != nil {
		logs.Error("添加参数执行出错,错误信息：", err)
		resAttr.Meta = &models.ResMeta{"添加参数执行出错", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 返回成功信息
	resAttr.Data = attr
	resAttr.Meta = &models.ResMeta{"添加参数成功", 201}
	this.Data["json"] = resAttr
	this.ServeJSON()
}

// 修改动态参数或者静态属性 接口：categories/:id/attributes/:attrId 请求方式：put
func (this *GoodsController) UpdateAttr() {
	var resAttr models.ResAttr

	// 没有找到修改参数相应的权限，此处省略权限验证

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		logs.Error("分类id格式错误")
		resAttr.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	if !models.CateIdExist(id) {
		logs.Error("分类id不存在")
		resAttr.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	attrId, err := this.GetInt(":attrId")
	if err != nil {
		logs.Error("属性id格式错误")
		resAttr.Meta = &models.ResMeta{"属性id格式错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 检验当前分类id下的属性id是否存在
	if !models.AttrIdExist(id, attrId) {
		logs.Error("属性id不存在")
		resAttr.Meta = &models.ResMeta{"属性id不存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	var addAttrParams models.AddAttrParams
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &addAttrParams)
	if err != nil {
		logs.Error("参数解析错误")
		resAttr.Meta = &models.ResMeta{"参数解析错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 验证请求体中的参数是否包含Attr_vals字段
	var validParamExist models.ValidParamExist
	var valsInBody bool
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &validParamExist)
	if err != nil {
		logs.Error(err)
		valsInBody = true
	} else {
		valsInBody = false
	}

	if addAttrParams.Attr_name == "" {
		logs.Error("新属性的名称不能为空")
		resAttr.Meta = &models.ResMeta{"新属性的名称不能为空", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	if addAttrParams.Attr_sel == "" {
		logs.Error("属性类型不能为空")
		resAttr.Meta = &models.ResMeta{"属性类型不能为空", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	} else if addAttrParams.Attr_sel != "only" && addAttrParams.Attr_sel != "many" {
		logs.Error("属性类型必须是'only'或者'many'")
		resAttr.Meta = &models.ResMeta{"属性类型必须是'only'或者'many'", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 修改参数时验证某一分类下的分类名称是否存在
	if models.UpdateAttrExist(id, attrId, addAttrParams.Attr_name, addAttrParams.Attr_sel) {
		logs.Error("属性名称已经存在")
		resAttr.Meta = &models.ResMeta{"属性名称已经存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 修改参数时验证属性类型是否正确
	if !models.AttrSelExist(attrId, addAttrParams.Attr_sel) {
		logs.Error("属性类型错误")
		resAttr.Meta = &models.ResMeta{"属性类型错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 数据验证完毕，执行数据库修改操作
	attr, err := models.UpdateAttr(attrId, addAttrParams.Attr_name, addAttrParams.Attr_vals, valsInBody)
	if err != nil {
		logs.Error("修改参数执行出错,错误信息：", err)
		resAttr.Meta = &models.ResMeta{"修改参数执行出错", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 返回成功信息
	resAttr.Data = attr
	resAttr.Meta = &models.ResMeta{"修改参数成功", 200}
	this.Data["json"] = resAttr
	this.ServeJSON()
}

// 删除动态参数或者静态属性 接口：categories/:id/attributes/:attrId 请求方式：delete
func (this *GoodsController) DeleteAttr() {
	var resAttr models.ResAttr

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 157)
	if !hasRight {
		logs.Error("权限不足")
		resAttr.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		logs.Error("分类id格式错误")
		resAttr.Meta = &models.ResMeta{"分类id格式错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 检验分类id是否存在
	if !models.CateIdExist(id) {
		logs.Error("分类id不存在")
		resAttr.Meta = &models.ResMeta{"分类id不存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	attrId, err := this.GetInt(":attrId")
	if err != nil {
		logs.Error("属性id格式错误")
		resAttr.Meta = &models.ResMeta{"属性id格式错误", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}
	// 检验当前分类id下的属性id是否存在
	if !models.AttrIdExist(id, attrId) {
		logs.Error("属性id不存在")
		resAttr.Meta = &models.ResMeta{"属性id不存在", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 数据验证完毕，执行数据库删除操作
	err = models.DeleteAttr(attrId)
	if err != nil {
		logs.Error("删除参数执行出错,错误信息：", err)
		resAttr.Meta = &models.ResMeta{"删除参数执行出错", 400}
		this.Data["json"] = resAttr
		this.ServeJSON()
		return
	}

	// 返回成功信息
	resAttr.Meta = &models.ResMeta{"删除参数成功", 200}
	this.Data["json"] = resAttr
	this.ServeJSON()
}

// 获取商品列表 接口：goods 请求方式：get
func (this *GoodsController) GetGoodsList() {
	var resGoodsList models.ResGoodsList

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 153)
	if !hasRight {
		logs.Error("权限不足")
		resGoodsList.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}

	/* 获取数据 */
	query := this.GetString("query")
	pagenum, err := this.GetInt("pagenum")
	if err != nil {
		logs.Error("pagenum为空或类型错误")
		resGoodsList.Meta = &models.ResMeta{"pagenum为空或类型错误", 400}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}
	pagesize, err := this.GetInt("pagesize")
	if err != nil {
		logs.Error("pagesize为空或类型错误")
		resGoodsList.Meta = &models.ResMeta{"pagesize为空或类型错误", 400}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}
	if pagenum <= 0 {
		logs.Error("pagenum必须大于0")
		resGoodsList.Meta = &models.ResMeta{"pagenum必须大于0", 400}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}
	if pagesize <= 0 {
		logs.Error("pagesize必须大于0")
		resGoodsList.Meta = &models.ResMeta{"pagesize必须大于0", 400}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}

	// 验证完毕，执行查询操作
	total, goodsList, err := models.GetGoodsList(query, pagenum, pagesize)
	if err != nil {
		logs.Error(err)
		resGoodsList.Meta = &models.ResMeta{"查询执行出错", 400}
		this.Data["json"] = resGoodsList
		this.ServeJSON()
		return
	}

	resGoodsList.Data = &models.ResGoodsData{total, pagenum, goodsList}
	resGoodsList.Meta = &models.ResMeta{"获取商品数据列表成功", 200}
	this.Data["json"] = resGoodsList
	this.ServeJSON()
}

// 上传图片 接口：【请求路径：upload 请求方式：post】
func (this *GoodsController) UploadPicture() {
	var resUpload models.ResUpload

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 150)
	if !hasRight {
		logs.Error("权限不足")
		resUpload.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resUpload
		this.ServeJSON()
		return
	}

	oriPicPath, err := UploadFile(&this.Controller, "uploadPicture")
	if err != nil {
		resUpload.Meta = &models.ResMeta{err.Error(), 400}
		this.Data["json"] = resUpload
		this.ServeJSON()
		return
	}
	resUpload.Data = &models.ResUploadData{oriPicPath, "http://127.0.0.1:8700" + oriPicPath}
	resUpload.Meta = &models.ResMeta{"上传图片成功", 200}
	this.Data["json"] = resUpload
	this.ServeJSON()

	// 以下代码应该在添加商品接口中使用，此处仅用作测试
	////读取本地文件
	//imgData,err:=ioutil.ReadFile("."+oriPicPath)
	//if err!=nil{
	//	logs.Error("图片读取错误：",err)
	//}
	//buf:=bytes.NewBuffer(imgData)
	//image,err:=imaging.Decode(buf)
	//if err!=nil{
	//	logs.Error("图片解码错误",err)
	//}
	////生成缩略图，尺寸800*800，并保存文件
	//image=imaging.Resize(image, 800, 800, imaging.Lanczos)
	//largePicPath:=strings.Replace(oriPicPath,".","_800×800.",-1)
	//err=imaging.Save(image,"."+largePicPath)
	//if err!=nil {
	//	logs.Error("图片保存错误",err)
	//}
	//this.Ctx.WriteString("success")
}

// 封装上传文件函数
func UploadFile(this *beego.Controller, filePath string) (string, error) {
	// GetFile函数返回的三个对象分别是：文件字节流、文件头（结构体：包含文件大小、文件名称等信息）、错误信息
	file, head, err := this.GetFile(filePath)
	if head == nil {
		return "", errors.New("获取文件头时出错")
	}
	if err != nil {
		logs.Error("读取文件时出错，错误原因：", err)
		return "", errors.New("读取文件时出错")
	}
	defer file.Close()
	logs.Info(fmt.Sprintf("文件大小：%vB,%vK,%vM", head.Size, head.Size/1024, head.Size/1024/1024))

	// 一般在后台执行文件上传之前还需要做一些相应的处理：
	// 01、判断文件大小，如果文件太大则阻止上传操作。这里文件大小不能超过 500 kb
	if head.Size > 500*1024 {
		logs.Error("文件太大，请重新上传")
		return "", errors.New("文件太大，请重新上传")
	}

	// 02、判断文件格式，这里必须是图片格式
	// 获取字符串中的后缀名
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		logs.Error("文件格式错误，请重新上传")
		return "", errors.New("文件格式错误，请重新上传")
	}

	// 03、防止重名，因为很多用户都上传文件很容易导致重名的情况
	fileName := xid.New().String()

	// 原始图片的存储路径
	oriPicPath := "/static/img/" + fileName + ext

	// SaveToFile这个函数用于存储用户上传过来的文件（实体），数据库中存储的是文件路径
	// 第一个参数是前端页面中上传文件组件的name属性值，第二个参数是存储路径
	// 第二个参数的指定有些怪异，可能是beego框架的一个小bug：存储文件的时候要在路径前加一个点，取的时候就不能有这个点
	err = this.SaveToFile(filePath, "."+oriPicPath)
	if err != nil {
		logs.Error("存储文件错误：", err)
		return "", errors.New("存储文件错误：" + err.Error())
	}
	return oriPicPath, nil
}

// 添加商品 接口：【请求路径：goods 请求方式：post】
func (this *GoodsController) AddGood() {
	var resAddGood models.ResAddGood

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 105)
	if !hasRight {
		logs.Error("权限不足")
		resAddGood.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resAddGood
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	var addGoodBody models.AddGoodBody
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &addGoodBody)
	if err != nil {
		logs.Error("参数解析错误:", err.Error())
		resAddGood.Meta = &models.ResMeta{"参数解析错误", 400}
		this.Data["json"] = resAddGood
		this.ServeJSON()
		return
	}

	resAddGoodData, err := models.AddGood(addGoodBody)
	if err != nil {
		logs.Error(err)
		resAddGood.Meta = &models.ResMeta{"添加商品失败", 400}
		this.Data["json"] = resAddGood
		this.ServeJSON()
		return
	}

	resAddGood.Data = resAddGoodData
	resAddGood.Meta = &models.ResMeta{"添加商品成功", 200}
	this.Data["json"] = resAddGood
	this.ServeJSON()
}

// 修改商品信息
func (this *GoodsController) UpdateGoodInfo() {
	var resGoodInfo models.ResGoodInfo
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 116)
	if !hasRight {
		resGoodInfo.Meta = &models.ResMeta{"权限不足", 403}
		logs.Error("权限不足")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resGoodInfo.Meta = &models.ResMeta{"商品id错误", 400}
		logs.Error("商品id错误")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 验证商品id是否存在
	goodExists := models.GoodExists(id)
	if !goodExists {
		resGoodInfo.Meta = &models.ResMeta{"商品id不存在", 400}
		logs.Error("商品id不存在")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	//good:=models.SpGoods{}
	good := new(models.SpGoods)
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &good)
	if err != nil {
		resGoodInfo.Meta = &models.ResMeta{"请求体中参数错误", 400}
		logs.Error("请求体中的参数错误")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}
	if good.GoodsName == "" {
		resGoodInfo.Meta = &models.ResMeta{"商品名称不能为空", 400}
		logs.Error("商品名称不能为空")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 执行数据库更新操作
	resGood, err := models.UpdateGood(good, id)
	if err != nil {
		resGoodInfo.Meta = &models.ResMeta{"修改商品失败", 400}
		logs.Error("修改商品失败", err)
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}
	resGoodInfo.Data = &models.ResGoodInfoData{resGood.GoodsName, resGood.GoodsPrice, resGood.GoodsWeight}
	resGoodInfo.Meta = &models.ResMeta{"修改商品成功", 200}
	this.Data["json"] = resGoodInfo
	this.ServeJSON()
	return
}

// 删除商品接口
func (this *GoodsController) DeleteGood() {
	var resGoodInfo models.ResGoodInfo
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 117)
	if !hasRight {
		resGoodInfo.Meta = &models.ResMeta{"权限不足", 403}
		logs.Error("权限不足")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resGoodInfo.Meta = &models.ResMeta{"商品id错误", 400}
		logs.Error("商品id错误")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 验证商品id是否存在
	goodExists := models.GoodExists(id)
	if !goodExists {
		resGoodInfo.Meta = &models.ResMeta{"商品id不存在", 400}
		logs.Error("商品id不存在")
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}

	// 执行数据库更新操作（假删）
	err = models.DeleteGood(id)
	if err != nil {
		resGoodInfo.Meta = &models.ResMeta{"删除商品失败", 400}
		logs.Error("修改商品失败", err)
		this.Data["json"] = resGoodInfo
		this.ServeJSON()
		return
	}
	resGoodInfo.Meta = &models.ResMeta{"删除商品成功", 200}
	this.Data["json"] = resGoodInfo
	this.ServeJSON()
	return
}
