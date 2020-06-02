package controllers

import (
	"JDStore/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
)

type RightsController struct {
	beego.Controller
}

// 获取权限列表
func (this *RightsController) GetRightsList() {
	var validateRight models.ResRoles
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 112)
	if !hasRight {
		validateRight.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = validateRight
		this.ServeJSON()
		return
	}
	// 获取数据
	GetRightsType := this.GetString(":type")

	// 查询数据
	switch GetRightsType {
	case "list":
		var resListRights models.ResListRights
		rightsList, err := models.QueryRightsList()
		if err != nil {
			resListRights.Meta = &models.ResUsersMeta{"查询执行错误", 400}
			this.Data["json"] = resListRights
			this.ServeJSON()
			return
		}
		resListRights.Data = &rightsList
		resListRights.Meta = &models.ResUsersMeta{"获取权限列表成功", 200}
		this.Data["json"] = resListRights
		this.ServeJSON()
		return
	case "tree":
		var resTreeRights models.ResTreeRights
		rightsTree, err := models.RightsTree()
		if err != nil {
			resTreeRights.Meta = &models.ResUsersMeta{"查询执行错误", 400}
			this.Data["json"] = resTreeRights
			this.ServeJSON()
			return
		}
		resTreeRights.Data = rightsTree
		resTreeRights.Meta = &models.ResUsersMeta{"获取权限列表成功", 200}
		this.Data["json"] = resTreeRights
		this.ServeJSON()
		return
	default:
		var resListRights models.ResListRights
		resListRights.Meta = &models.ResUsersMeta{"显示类型参数错误", 400}
		this.Data["json"] = resListRights
		this.ServeJSON()
	}
}

// 获取角色列表
func (this *RightsController) GetRolesList() {
	var resRoles models.ResRoles
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 138)
	if !hasRight {
		resRoles.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}
	roleList, err := models.QueryRoleList()
	if err != nil {
		logs.Error(err)
		resRoles.Meta = &models.ResUsersMeta{"查询执行失败", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}
	resRoles.Data = roleList
	resRoles.Meta = &models.ResUsersMeta{"获取角色列表成功", 200}
	this.Data["json"] = resRoles
	this.ServeJSON()
}

// 删除角色指定权限
func (this *RightsController) DeleteRight() {
	var resDelRight models.ResDelRight
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 135)
	if !hasRight {
		resDelRight.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}
	// 获取数据
	roleId, err := this.GetInt(":roleId")
	if err != nil {
		resDelRight.Meta = &models.ResUsersMeta{"参数roleId错误", 400}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}
	rightId, err := this.GetInt(":rightId")
	if err != nil {
		resDelRight.Meta = &models.ResUsersMeta{"参数rightId错误", 400}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}

	// 数据处理
	// 查询该角色信息
	spRole := models.SpRole{RoleId: roleId}
	o := orm.NewOrm()
	err = o.Read(&spRole, "RoleId")
	if err != nil {
		resDelRight.Meta = &models.ResUsersMeta{"角色不存在", 400}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}
	// 判断参数中的权限id是否存在于该角色的所有权限之中
	isInRights := false
	roleRights := strings.Split(spRole.PsIds, ",")
	// 记录权限对应的索引，便于后续删除
	var index int
	for i, v := range roleRights {
		if v == strconv.FormatInt(int64(rightId), 10) {
			index = i
			isInRights = true
			break
		}
	}
	// 没有在角色列表中找到该权限
	if !isInRights {
		resDelRight.Meta = &models.ResUsersMeta{"权限不存在", 400}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}
	// 权限存在，执行删除
	roleRights = append(roleRights[:index], roleRights[index+1:]...)
	// 更新数据库
	spRole.PsIds = strings.Join(roleRights, ",")
	_, err = o.Update(&spRole, "PsIds")
	if err != nil {
		resDelRight.Meta = &models.ResUsersMeta{"更新数据库执行出错", 400}
		this.Data["json"] = resDelRight
		this.ServeJSON()
		return
	}
	// 更新成功，返回该角色的权限列表
	resRoleRights, err := models.GetRoleRights(0, roleRights)
	resDelRight.Data = resRoleRights
	resDelRight.Meta = &models.ResUsersMeta{"删除权限成功", 200}
	this.Data["json"] = resDelRight
	this.ServeJSON()
}

// 更新某个角色的权限
func (this *RightsController) UpdateRoleRights() {
	var resRoles models.ResRoles
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 141)
	if !hasRight {
		resRoles.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	// 获取数据
	// 获取角色id
	roleId, err := this.GetInt(":roleId")
	if err != nil {
		resRoles.Meta = &models.ResUsersMeta{"角色id参数错误", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	//校验请求参数中的角色id是否存在于角色列表
	var spRoles []models.SpRole
	o := orm.NewOrm()
	_, err = o.QueryTable("sp_role").All(&spRoles)
	if err != nil {
		resRoles.Meta = &models.ResUsersMeta{"查询执行出错", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}
	isNotInTable := true
	for _, v := range spRoles {
		if v.RoleId == roleId {
			isNotInTable = false
			break
		}
	}
	if isNotInTable {
		resRoles.Meta = &models.ResUsersMeta{"角色不存在", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	// 获取权限列表
	var postRights models.PostRights
	postData := this.Ctx.Input.RequestBody
	err = json.Unmarshal(postData, &postRights)
	if err != nil {
		resRoles.Meta = &models.ResUsersMeta{"json解析错误", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	// 校验请求参数中的权限列表
	valid := validation.Validation{}
	b, err := valid.Valid(&postRights)
	if err != nil {
		// handle error
		resRoles.Meta = &models.ResUsersMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}
	if !b {
		// validation does not pass
		// blabla...
		msg := ""
		// 出错就退出循环，不必再继续遍历
		for _, err := range valid.Errors {
			//logs.Error(err.Field, err.Name, err.Key, err.Message)
			// err.Field是出错字段，err.Name是校验方法，err.Key是出错字段+校验方法，err.Message是错误信息
			msg = fmt.Sprintf("错误字段：%s，错误信息：%s", err.Field, err.Message)
			break
		}
		resRoles.Meta = &models.ResUsersMeta{msg, 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	// 请求参数验证通过，更新数据库
	var role models.SpRole
	role.RoleId = roleId
	err = o.Read(&role)
	if err != nil {
		resRoles.Meta = &models.ResUsersMeta{"查询角色出错", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}
	/* 赋予根权限0 */
	role.PsIds = "0," + postRights.Rids
	_, err = o.Update(&role, "PsIds")
	if err != nil {
		resRoles.Meta = &models.ResUsersMeta{"更新字段值出错", 400}
		this.Data["json"] = resRoles
		this.ServeJSON()
		return
	}

	// 返回更新成功信息
	resRoles.Meta = &models.ResUsersMeta{"更新角色权限成功", 200}
	this.Data["json"] = resRoles
	this.ServeJSON()
}

// 修改角色信息
func (this *RightsController) UpdateRoleInfo() {
	var resRoleInfo models.ResRoleInfo
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 140)
	if !hasRight {
		resRoleInfo.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}

	//获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"角色id错误", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	var roleParams models.GetRoleParams
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &roleParams)
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"请求体中参数错误", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	if roleParams.RoleName == "" {
		resRoleInfo.Meta = &models.ResMeta{"角色名称不能为空", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}

	// 更新数据库
	// 查询角色信息
	role := models.SpRole{RoleId: id}
	o := orm.NewOrm()
	err = o.Read(&role, "RoleId")
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"角色id不存在", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 更新角色信息
	role.RoleName = roleParams.RoleName
	role.RoleDesc = roleParams.RoleDesc
	_, err = o.Update(&role, "RoleName", "RoleDesc")
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"更新执行出错", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 返回成功信息
	resRoleInfo.Data = &models.ResRoleInfoData{role.RoleId, role.RoleName, role.RoleDesc}
	resRoleInfo.Meta = &models.ResMeta{"更新角色信息成功", 200}
	this.Data["json"] = resRoleInfo
	this.ServeJSON()
}

// 删除角色信息
func (this *RightsController) DeleteRole() {
	var resRoleInfo models.ResRoleInfo
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 130)
	if !hasRight {
		resRoleInfo.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"角色id错误", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 删除角色
	role := models.SpRole{RoleId: id}
	o := orm.NewOrm()
	num, err := o.Delete(&role, "RoleId")
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"删除执行失败", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	if num == 0 {
		resRoleInfo.Meta = &models.ResMeta{"角色id不存在", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}

	// 返回成功信息
	resRoleInfo.Meta = &models.ResMeta{"删除角色成功", 200}
	this.Data["json"] = resRoleInfo
	this.ServeJSON()
}

// 添加角色
func (this *RightsController) AddRole() {
	var resRoleInfo models.ResRoleInfo
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 129)
	if !hasRight {
		resRoleInfo.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}

	// 获取数据
	var roleParams models.GetRoleParams
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &roleParams)
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"请求体中参数错误", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	if roleParams.RoleName == "" {
		resRoleInfo.Meta = &models.ResMeta{"角色名称不能为空", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 验证要添加的角色名称是否已存在
	role := models.SpRole{RoleName: roleParams.RoleName, RoleDesc: roleParams.RoleDesc}
	o := orm.NewOrm()
	err = o.Read(&role, "RoleName")
	if err == nil {
		resRoleInfo.Meta = &models.ResMeta{"角色名称已存在", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 添加角色到数据库
	_, err = o.Insert(&role)
	if err != nil {
		resRoleInfo.Meta = &models.ResMeta{"添加执行错误", 400}
		this.Data["json"] = resRoleInfo
		this.ServeJSON()
		return
	}
	// 返回成功信息
	resRoleInfo.Data = &models.ResRoleInfoData{role.RoleId, role.RoleName, role.RoleDesc}
	resRoleInfo.Meta = &models.ResMeta{"添加角色成功", 200}
	this.Data["json"] = resRoleInfo
	this.ServeJSON()
}
