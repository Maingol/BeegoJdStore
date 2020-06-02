package controllers

import (
	"JDStore/models"
	"github.com/astaxie/beego"
)

type MenusController struct {
	beego.Controller
}

func (this *MenusController) HandleGetMenus() {
	var resMenus models.ResMenus
	menuLv1, err := models.TestOneToOne(1, 0)
	if err != nil {
		resMenus.Meta = &models.ResMeta{"查询数据库时出现错误", 500}
		return
	}
	/* 不可以通过遍历对象的value来改变这个value的值，
	*  一定要通过遍历对象的索引来改变value的值，
	*  因为range遍历的是对象的副本，而不是对象的引用 */
	for i := range menuLv1 {
		menuLv1[i].Children, err = models.TestOneToOne(2, menuLv1[i].PsId)
		if err != nil {
			resMenus.Meta = &models.ResMeta{"查询数据库时出现错误", 500}
			return
		}
	}
	resMenus.Data = menuLv1
	resMenus.Meta = &models.ResMeta{"获取菜单列表成功", 200}
	this.Data["json"] = resMenus
	this.ServeJSON()
}
