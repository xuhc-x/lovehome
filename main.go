package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "lovehome/models"
	_ "lovehome/routers"
	"net/http"
	"strings"
)

func main() {
	ignoreStaticPath()
	//models.UpLoadFile("main.go")
	//models.TestUploadByFilename("main.go")
	beego.BConfig.WebConfig.Session.SessionOn = true  //使用session模块的配置，设置开启
	beego.Run(":8098")
}

func ignoreStaticPath() {

	//透明static
	beego.SetStaticPath("group1/M00/","fdfs/storage_data/data/")

	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)  //分离筛选，直接跳到页面
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//fmt.Print("request url: ", orrpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}


