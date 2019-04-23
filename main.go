package main

import (
	"os"

	"upgrade.maiyajia.com/controllers"
	_ "upgrade.maiyajia.com/routers"
	"upgrade.maiyajia.com/services/mongo"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		//设置应用日志
		beego.SetLogger(logs.AdapterFile, `{"filename":"app.log"}`)
		beego.SetLevel(beego.LevelInformational)
		// beego.BConfig.WebConfig.DirectoryIndex = true
		// beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 设置跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	//启动mongodb数据库，初始化mongo数据库连接会话
	if err := mongo.Startup(); err != nil {
		beego.Error(err)
		os.Exit(0)
	}
	beego.Informational("数据库启动成功")

	// 注册错误处理函数
	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}
