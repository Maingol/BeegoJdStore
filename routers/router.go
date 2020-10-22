package routers

import (
	"JDStore/controllers"
	"JDStore/models"
	"JDStore/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"strings"
)

func init() {
	// 获取配置文件中定义的基准url
	baseURL := beego.AppConfig.String("baseURL")

	// 实现跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		/*这里最后一个元素“mytoken”是想在请求头中自定义的key，发请求时所有在请求头中自定义的key都必须写在这里*/
		//AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type","mytoken"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	// 路由拦截器
	beego.InsertFilter(baseURL+"*", beego.BeforeRouter, func(ctx *context.Context) {
		// jwt鉴权
		var token string
		if strings.Index(ctx.Request.RequestURI, "/login") >= 0 {
			// url中包含/login  ，放行
			return
		} else {
			// 获取请求头中的token
			token = ctx.Input.Header("Authorization")
			// 去掉token值前面的“Bearer ”部分，截取出后面真正的token值
			if token != "" {
				token = strings.Split(token, " ")[1]
				// logs.Info("token:",token)
			}
			logs.Info("current router path is ", ctx.Request.RequestURI)
		}
		if token == "" || utils.CheckToken(token) == "" {
			// 验证未通过（请求头中不含有令牌，或者令牌未通过验证）
			var resLogin models.ResLogin
			resLogin.Meta = &models.ResMeta{"无效token", 401}
			logs.Error("请求中不包含token或token无效")
			res, _ := json.Marshal(resLogin)
			ctx.ResponseWriter.Write(res)
		} else if utils.CheckToken(token) != "" {
			// 验证token通过
			logs.Info(fmt.Sprintf("用户:%v 鉴权成功", utils.CheckToken(token)))
		}
	})

	beego.Router("/", &controllers.MainController{})
	beego.Router(baseURL+"login", &controllers.LoginController{}, "post:HandlePost")
	beego.Router(baseURL+"menus", &controllers.MenusController{}, "get:HandleGetMenus")
	beego.Router(baseURL+"users", &controllers.UsersController{},
		"get:HandleGetUsers;post:AddUser")
	beego.Router(baseURL+"users/:uId/state/:type", &controllers.UsersController{}, "put:PutUserState")
	beego.Router(baseURL+"users/:id", &controllers.UsersController{},
		"get:GetUserInfo;put:UpdateUserInfo;delete:DeleteUser")
	beego.Router(baseURL+"rights/:type", &controllers.RightsController{}, "get:GetRightsList")
	beego.Router(baseURL+"roles", &controllers.RightsController{},
		"get:GetRolesList;post:AddRole")
	beego.Router(baseURL+"roles/:roleId/rights/:rightId", &controllers.RightsController{},
		"delete:DeleteRight")
	beego.Router(baseURL+"roles/:roleId/rights", &controllers.RightsController{},
		"post:UpdateRoleRights")
	beego.Router(baseURL+"users/:id/role", &controllers.UsersController{}, "put:UpdateRole")
	beego.Router(baseURL+"roles/:id", &controllers.RightsController{},
		"put:UpdateRoleInfo;delete:DeleteRole")
	beego.Router(baseURL+"categories", &controllers.GoodsController{},
		"get:GetGoodsCate;post:AddGoodsCate")
	beego.Router(baseURL+"categories/:id", &controllers.GoodsController{},
		"put:UpdateCateName;delete:DeleteCate")
	beego.Router(baseURL+"categories/:id/attributes", &controllers.GoodsController{},
		"get:GetAttrList;post:AddAttr")
	beego.Router(baseURL+"categories/:id/attributes/:attrId", &controllers.GoodsController{},
		"put:UpdateAttr;delete:DeleteAttr")
	beego.Router(baseURL+"goods", &controllers.GoodsController{},
		"get:GetGoodsList;post:AddGood")
	beego.Router(baseURL+"upload", &controllers.GoodsController{},
		"post:UploadPicture")
	beego.Router(baseURL+"goods/:id", &controllers.GoodsController{},
		"put:UpdateGoodInfo;delete:DeleteGood")
}
