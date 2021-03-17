package main

import (
	"BeeDemo/models"
	_ "BeeDemo/models" //不引入无法初始化数据库模型（创建表）
	_ "BeeDemo/routers"
	"BeeDemo/utils"
	"encoding/gob"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	//开启sql查询日志
	orm.Debug = true
	//session 设置redis驱动
	gob.Register(&models.User{})
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "users"
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "127.0.0.1:6379"

	//页面登录过滤器，检测用户是否登录
	beego.InsertFilter("/cms/system/*", beego.BeforeRouter, utils.CmsLoginFilter)
	orm.RunCommand()
	beego.Run()
}

func init() {
	MySQL()
	utils.Redis()
}

//初始化 MySQL
func MySQL() {
	name, _ := beego.AppConfig.String("username")
	pwd, _ := beego.AppConfig.String("password")
	host, _ := beego.AppConfig.String("host")
	port, _ := beego.AppConfig.String("port")
	dbname, _ := beego.AppConfig.String("dbname")

	err := orm.RegisterDriver("mysql", orm.DRMySQL)

	if err != nil {
		fmt.Println("数据库驱动注册失败：", err)
		return
	}

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", name, pwd, host, port, dbname)
	err = orm.RegisterDataBase("default", "mysql", source)
	if err != nil {
		fmt.Println("数据库链接失败：", err)
		return
	}

	//同步数据库模型表结构
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println("同步数据库表结构失败：", err)
		return
	}

	println("数据库初始化完成...")
}
