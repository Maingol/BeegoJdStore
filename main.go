package main

import (
	_ "JDStore/routers"
	"JDStore/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbHost := beego.AppConfig.String("dbHost")
	dbPort := beego.AppConfig.String("dbPort")
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")
	dbName := beego.AppConfig.String("dbName")

	//注册mysql Driver
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		logs.Error("注册mysql驱动失败", err)
	}
	//构造conn连接
	//用户名:密码@tcp(url地址)/数据库
	conn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	//注册数据库连接
	err = orm.RegisterDataBase("default", "mysql", conn)
	if err != nil {
		logs.Error("注册mysql连接失败", err)
	}
}

func main() {
	// 重新定义表单校验的错误信息
	utils.ErrorMessage()

	beego.Run()
}
