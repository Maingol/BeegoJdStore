package controllers

import (
	"JDStore/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ReportController struct {
	beego.Controller
}

// 获取数据报表 【接口：reports/type/1  请求方式：get】
func (this ReportController) GetReport() {
	var resReport models.ResReport
	// 权限验证
	hasRight := models.ValidateRight(this.Ctx, 145)
	if !hasRight {
		resReport.Meta = &models.ResMeta{"权限不足", 403}
		logs.Error("权限不足")
		this.Data["json"] = resReport
		this.ServeJSON()
		return
	}

	resReportData, err := models.GetReport()
	if err != nil {
		logs.Error("获取数据报表失败", err)
		resReport.Meta = &models.ResMeta{"获取数据报表失败", 400}
		this.Data["json"] = resReport
		this.ServeJSON()
		return
	}

	resReport.Data = resReportData
	resReport.Meta = &models.ResMeta{"获取数据报表成功", 200}
	this.Data["json"] = resReport
	this.ServeJSON()
	return
}
