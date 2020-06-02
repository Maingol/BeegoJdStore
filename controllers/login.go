package controllers

import (
	"JDStore/models"
	"JDStore/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) HandlePost() {
	/*获取post请求中的参数（非json格式）*/
	//var user GetParams
	//user.Username = this.GetString("username")
	//user.Password = this.GetString("password")

	/*获取post请求中的参数（json格式）*/
	var user models.GetParams
	postData := this.Ctx.Input.RequestBody
	err := json.Unmarshal(postData, &user)
	if err != nil {
		beego.Info("json.Unmarshal is err:", err.Error())
	}
	//logs.Info(user)

	/*数据校验*/
	valid := validation.Validation{}
	//自定义消息提示，官方默认的是英文，我们手动改成中文,具体内容可以查看 MessageTmpls 。
	//var MessageTmpls = map[string]string{
	//	"Required": "不能为空",
	//	"MinSize":  "最短长度为 %d",
	//	"Length":   "长度必须为 %d",
	//	"Numeric":  "必须是有效的数字",
	//	"Email":    "必须是有效的电子邮件地址",
	//	"Mobile":   "必须是有效的手机号码",
	//}
	//validation.SetDefaultMessage(MessageTmpls)
	b, err := valid.Valid(&user)
	var resLogin models.ResLogin
	if err != nil {
		// handle error
		resLogin.Meta = &models.ResMeta{"数据校验时发生内部错误", 500}
		this.Data["json"] = resLogin
		this.ServeJSON()
		return
	}
	if !b {
		// validation does not pass
		// blabla...
		for _, err := range valid.Errors {
			logs.Error(err.Key, err.Message)
		}
		resLogin.Meta = &models.ResMeta{"用户名或密码格式错误", 422}
		this.Data["json"] = resLogin
		this.ServeJSON()
		return
	}

	/*从数据库中查询用户名是否存在以及密码是否正确*/
	o := orm.NewOrm()
	manager := models.SpManager{MgName: user.UserName}
	err = o.Read(&manager, "MgName")
	if err != nil {
		resLogin.Meta = &models.ResMeta{"用户不存在", 400}
		this.Data["json"] = resLogin
		this.ServeJSON()
		return
	}
	/*将请求参数中的密码和数据库中加密后的密码进行比较，这里用到的是封装好的函数。
	* 关于如何使用hash对密码进行加密、解密的详情参见：https://www.kancloud.cn/golang_programe/golang/1144844*/
	data := []byte(user.PassWord)
	result := utils.ComparePasswords(manager.MgPwd, data)
	if !result {
		resLogin.Meta = &models.ResMeta{"密码错误", 400}
		this.Data["json"] = resLogin
		this.ServeJSON()
		return
	} else {
		/*程序执行到这里说明，用户存在且密码正确。创建一个token，将其和用户信息一同返回给前端，并返回成功信息。*/
		token := utils.CreateToken(manager.MgName)
		resLogin.Data = &models.ResData{manager.MgId, manager.RoleId, manager.MgName,
			manager.MgMobile, manager.MgEmail, utils.GetHeaderTokenValue(token)}
		resLogin.Meta = &models.ResMeta{"登陆成功", 200}
		this.Data["json"] = resLogin
		this.ServeJSON()
		return
	}
}
