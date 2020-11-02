package controllers

import (
	"github.com/astaxie/beego"
)

// 这是用户访问首页的控制器

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	// 访问首页
	c.TplName = "index.html"
}
