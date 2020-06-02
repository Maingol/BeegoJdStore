package controllers

import (
	"JDStore/models"
	"JDStore/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strings"
	"time"
)

type UsersController struct {
	beego.Controller
}

// 获取用户列表
func (this *UsersController) HandleGetUsers() {
	var resUsers models.ResUsers
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 110)
	if !hasRight {
		resUsers.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}

	/* 获取数据 */
	query := this.GetString("query")
	pagenum, err := this.GetInt("pagenum")
	if err != nil {
		resUsers.Meta = models.ResUsersMeta{"pagenum 参数错误", 400}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}
	pagesize, err := this.GetInt("pagesize")
	if err != nil {
		resUsers.Meta = models.ResUsersMeta{"pagesize 参数错误", 400}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}
	/* 将参数赋值给定义好的结构体，准备校验 */
	params := models.UsersParams{pagenum, pagesize}

	/*数据校验*/
	valid := validation.Validation{}
	b, err := valid.Valid(&params)
	if err != nil {
		// handle error
		resUsers.Meta = models.ResUsersMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}
	if !b {
		// validation does not pass
		// blabla...
		for _, err := range valid.Errors {
			logs.Error(err.Key, err.Message)
			if strings.HasPrefix(err.Key, "Pagenum") {
				resUsers.Meta = models.ResUsersMeta{"pagenum 参数错误", 400}
				this.Data["json"] = resUsers
				this.ServeJSON()
				return
			} else if strings.HasPrefix(err.Key, "Pagesize") {
				resUsers.Meta = models.ResUsersMeta{"pagesize 参数错误", 400}
				this.Data["json"] = resUsers
				this.ServeJSON()
				return
			}
		}
	}

	// 查询管理员的总记录数
	o := orm.NewOrm()
	// total,err:=o.QueryTable("sp_manager").Count()
	// 替换上面的查询方式，使用模糊匹配
	total, err := o.QueryTable("sp_manager").Filter("mg_name__contains", query).Count()
	if err != nil {
		resUsers.Meta = models.ResUsersMeta{"查询总记录数出错", 400}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}
	// 查询管理员列表
	managers, err := models.GetManagers(query, pagenum, pagesize)
	if err != nil {
		resUsers.Meta = models.ResUsersMeta{"查询管理员列表出错", 400}
		this.Data["json"] = resUsers
		this.ServeJSON()
		return
	}
	resUsers.Data = models.ResUsersData{total, pagenum, managers}
	resUsers.Meta = models.ResUsersMeta{"获取管理员列表成功", 200}
	this.Data["json"] = resUsers
	this.ServeJSON()
}

// 修改用户的状态
func (this *UsersController) PutUserState() {
	var resUserState models.ResUserState
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 159)
	if !hasRight {
		resUserState.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resUserState
		this.ServeJSON()
		return
	}

	// 获取数据
	uId, err := this.GetInt(":uId")
	if err != nil {
		logs.Info(err)
		resUserState.Meta = &models.ResUsersMeta{"uId参数错误", 400}
		this.Data["json"] = resUserState
		this.ServeJSON()
	}

	state, err := this.GetBool(":type")
	if err != nil {
		logs.Info(err)
		resUserState.Meta = &models.ResUsersMeta{"type参数错误", 400}
		this.Data["json"] = resUserState
		this.ServeJSON()
	}

	// 这里无需做进一步数据校验，因为参数为空时，服务器直接报404错误

	o := orm.NewOrm()
	// 查询mg_id为uId的用户
	user := models.SpManager{MgId: uId}
	err = o.Read(&user)
	if err != nil {
		resUserState.Meta = &models.ResUsersMeta{"管理员ID不存在", 400}
		this.Data["json"] = resUserState
		this.ServeJSON()
	}

	// 修改字段值
	if state {
		user.MgState = 1
	} else {
		user.MgState = 0
	}

	// 执行更新操作
	_, err = o.Update(&user, "mg_state")
	if err != nil {
		resUserState.Meta = &models.ResUsersMeta{"修改执行出错", 400}
		this.Data["json"] = resUserState
		this.ServeJSON()
	}

	// 返回成功信息
	resUserState.Data = &models.ResStateData{
		user.MgId, user.RoleId, user.MgName,
		user.MgMobile, user.MgEmail, user.MgState}
	resUserState.Meta = &models.ResUsersMeta{"设置状态成功", 200}
	this.Data["json"] = resUserState
	this.ServeJSON()

}

// 添加新的用户
func (this *UsersController) AddUser() {
	var resAddUser models.ResAddUser
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 131)
	if !hasRight {
		resAddUser.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resAddUser
		this.ServeJSON()
		return
	}

	// 获取数据
	var user models.CreateUser
	postData := this.Ctx.Input.RequestBody
	err := json.Unmarshal(postData, &user)
	if err != nil {
		beego.Info("json.Unmarshal is err:", err.Error())
	}

	// 数据校验

	valid := validation.Validation{}

	b, err := valid.Valid(&user)
	if err != nil {
		// handle error
		resAddUser.Meta = &models.ResUsersMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resAddUser
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
		resAddUser.Meta = &models.ResUsersMeta{msg, 400}
		this.Data["json"] = resAddUser
		this.ServeJSON()
		return
	}

	// 检验用户名是否已存在
	manager := models.SpManager{MgName: user.UserName}
	o := orm.NewOrm()
	err = o.Read(&manager, "MgName")
	if err == nil {
		// 查询该用户没有出现异常，说明用户已存在
		resAddUser.Meta = &models.ResUsersMeta{"用户名已存在", 400}
		this.Data["json"] = resAddUser
		this.ServeJSON()
		return
	}
	// 用户名不存在，可以向数据库增加一条记录
	manager = models.SpManager{
		MgName: user.UserName, MgPwd: utils.HashAndSalt(user.PassWord),
		MgTime: int(time.Now().Unix()), MgMobile: user.Mobile, MgEmail: user.Email}
	_, err = o.Insert(&manager)
	if err != nil {
		resAddUser.Meta = &models.ResUsersMeta{"添加执行出错", 400}
		this.Data["json"] = resAddUser
		this.ServeJSON()
		return
	}

	// 返回新增用户的信息以及成功信息
	resAddUser.Data = &models.ResAddUserData{
		manager.MgId, manager.MgName, manager.MgMobile,
		manager.MgEmail, manager.RoleId, manager.MgTime}
	resAddUser.Meta = &models.ResUsersMeta{"创建成功", 201}
	this.Data["json"] = resAddUser
	this.ServeJSON()
}

// 根据id获取该用户信息
func (this *UsersController) GetUserInfo() {
	var resGetUser models.ResGetUser
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 136)
	if !hasRight {
		resGetUser.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":id")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"请检查id参数", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 查询数据
	user := models.SpManager{MgId: id}
	o := orm.NewOrm()
	err = o.Read(&user, "MgId")
	if err != nil {
		// 没有查到该用户，返回错误信息
		resGetUser.Meta = &models.ResUsersMeta{"用户不存在", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}
	// 查到了用户信息
	resGetUser.Data = &models.ResUserData{
		user.MgId, user.RoleId, user.MgName,
		user.MgMobile, user.MgEmail}
	resGetUser.Meta = &models.ResUsersMeta{"获取用户信息成功", 200}
	this.Data["json"] = resGetUser
	this.ServeJSON()

}

// 修改用户信息
func (this *UsersController) UpdateUserInfo() {
	var resGetUser models.ResGetUser
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 133)
	if !hasRight {
		resGetUser.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}
	// 获取数据
	// 获取用户id
	id, err := this.GetInt(":id")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"请检查id参数", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}
	// 获取请求参数
	var updateUser models.GetPutParams
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &updateUser)
	if err != nil {
		logs.Info("json.Unmarshal is err:", err)
	}

	// 变量的声明必须在goto之前
	var b bool
	valid := validation.Validation{}

	b, err = valid.Valid(&updateUser)
	if err != nil {
		// handle error
		resGetUser.Meta = &models.ResUsersMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resGetUser
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
		resGetUser.Meta = &models.ResUsersMeta{msg, 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 校验通过，修改用户信息
	manager := models.SpManager{MgId: id, MgEmail: updateUser.Email, MgMobile: updateUser.Mobile}
	o := orm.NewOrm()
	_, err = o.Update(&manager, "MgEmail", "MgMobile")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"修改用户信息失败", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 返回成功信息
	err = o.Read(&manager, "MgId")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"查询用户信息失败", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	resGetUser.Data = &models.ResUserData{
		manager.MgId, manager.RoleId, manager.MgName,
		manager.MgMobile, manager.MgEmail}
	resGetUser.Meta = &models.ResUsersMeta{"更新成功", 200}
	this.Data["json"] = resGetUser
	this.ServeJSON()
}

// 删除用户
func (this *UsersController) DeleteUser() {
	var resGetUser models.ResGetUser
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 132)
	if !hasRight {
		resGetUser.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}
	// 获取数据
	// 获取用户id
	id, err := this.GetInt(":id")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"请检查id参数", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 删除用户
	manager := models.SpManager{MgId: id}
	o := orm.NewOrm()
	num, err := o.Delete(&manager, "MgId")
	if err != nil {
		resGetUser.Meta = &models.ResUsersMeta{"删除执行错误", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}
	if num == 0 {
		resGetUser.Meta = &models.ResUsersMeta{"用户id不存在", 400}
		this.Data["json"] = resGetUser
		this.ServeJSON()
		return
	}

	// 删除成功
	resGetUser.Meta = &models.ResUsersMeta{"删除成功", 200}
	this.Data["json"] = resGetUser
	this.ServeJSON()
}

// 分配用户角色
func (this *UsersController) UpdateRole() {
	var resUser models.ResGetUser
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 134)
	if !hasRight {
		resUser.Meta = &models.ResUsersMeta{"权限不足", 403}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 获取数据
	userId, err := this.GetInt(":id")
	if err != nil {
		resUser.Meta = &models.ResUsersMeta{"管理员id格式错误", 400}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 检验用户id是否存在于数据库
	manager := models.SpManager{MgId: userId}
	o := orm.NewOrm()
	err = o.Read(&manager)
	if err != nil {
		resUser.Meta = &models.ResUsersMeta{"管理员id不存在", 400}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 获取请求体中的角色id
	var roleId models.RoleId
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &roleId)
	if err != nil {
		resUser.Meta = &models.ResUsersMeta{"请求体中不包含角色id或角色id不是整数", 400}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 校验请求参数中的角色id
	valid := validation.Validation{}
	b, err := valid.Valid(&roleId)
	if err != nil {
		// handle error
		resUser.Meta = &models.ResUsersMeta{"数据校验时发生内部错误", 400}
		this.Data["json"] = resUser
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
		resUser.Meta = &models.ResUsersMeta{msg, 400}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 校验成功，更新用户的角色
	manager.RoleId = int8(roleId.Rid)
	_, err = o.Update(&manager, "RoleId")
	if err != nil {
		resUser.Meta = &models.ResUsersMeta{"更新执行错误", 400}
		this.Data["json"] = resUser
		this.ServeJSON()
		return
	}

	// 返回成功信息
	resUser.Data = &models.ResUserData{manager.MgId, manager.RoleId,
		manager.MgName, manager.MgMobile, manager.MgEmail}
	resUser.Meta = &models.ResUsersMeta{"设置角色成功", 200}
	this.Data["json"] = resUser
	this.ServeJSON()
}
