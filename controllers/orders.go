package controllers

import (
	"JDStore/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type OrdersController struct {
	beego.Controller
}

// 获取订单列表 【接口：orders 请求方式：get】
func (this *OrdersController) GetOrdersList() {
	var resOrdersList models.ResOrdersList

	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 107)
	if !hasRight {
		logs.Error("权限不足")
		resOrdersList.Meta = &models.ResMeta{"权限不足", 403}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}

	/* 获取数据 */
	query := this.GetString("query")
	pagenum, err := this.GetInt("pagenum")
	if err != nil {
		logs.Error("pagenum为空或类型错误")
		resOrdersList.Meta = &models.ResMeta{"pagenum为空或类型错误", 400}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}
	pagesize, err := this.GetInt("pagesize")
	if err != nil {
		logs.Error("pagesize为空或类型错误")
		resOrdersList.Meta = &models.ResMeta{"pagesize为空或类型错误", 400}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}
	if pagenum <= 0 {
		logs.Error("pagenum必须大于0")
		resOrdersList.Meta = &models.ResMeta{"pagenum必须大于0", 400}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}
	if pagesize <= 0 {
		logs.Error("pagesize必须大于0")
		resOrdersList.Meta = &models.ResMeta{"pagesize必须大于0", 400}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}

	// 验证完毕，执行查询操作
	total, ordersList, err := models.GetOrdersList(query, pagenum, pagesize)
	if err != nil {
		logs.Error("查询执行出错", err)
		resOrdersList.Meta = &models.ResMeta{"查询执行出错", 400}
		this.Data["json"] = resOrdersList
		this.ServeJSON()
		return
	}

	resOrdersList.Data = &models.ResOrdersData{total, pagenum, ordersList}
	resOrdersList.Meta = &models.ResMeta{"获取订单数据列表成功", 200}
	this.Data["json"] = resOrdersList
	this.ServeJSON()
}

// 修改订单地址 【接口：orders:id 请求方式：put】
func (this *OrdersController) UpdateOrderAddr() {
	var resUpdateAddr models.ResUpdateAddr
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 154)
	if !hasRight {
		resUpdateAddr.Meta = &models.ResMeta{"权限不足", 403}
		logs.Error("权限不足")
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}

	// 获取数据
	id, err := this.GetInt(":order_id")
	if err != nil {
		resUpdateAddr.Meta = &models.ResMeta{"订单id错误", 400}
		logs.Error("订单id错误")
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}

	// 验证订单id是否存在
	orderExists := models.OrderExists(id)
	if !orderExists {
		resUpdateAddr.Meta = &models.ResMeta{"订单id不存在", 400}
		logs.Error("订单id不存在")
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}

	// 获取请求体中的参数
	order := new(models.SpOrder)
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &order)
	if err != nil {
		resUpdateAddr.Meta = &models.ResMeta{"请求体中参数错误", 400}
		logs.Error("请求体中的参数错误")
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}
	if order.ConsigneeAddr == "" {
		resUpdateAddr.Meta = &models.ResMeta{"订单地址不能为空", 400}
		logs.Error("订单地址不能为空")
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}

	// 执行数据库更新操作
	resAddr, err := models.UpdateAddr(order, id)
	if err != nil {
		resUpdateAddr.Meta = &models.ResMeta{"修改订单地址失败", 400}
		logs.Error("修改订单地址失败", err)
		this.Data["json"] = resUpdateAddr
		this.ServeJSON()
		return
	}
	resUpdateAddr.Data = resAddr
	resUpdateAddr.Meta = &models.ResMeta{"修改订单地址成功", 200}
	this.Data["json"] = resUpdateAddr
	this.ServeJSON()
	return
}
