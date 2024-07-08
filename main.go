package main

import (
	_ "github.com/udistrital/sga_plan_estudio_mid/routers"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	xray.InitXRay()
	beego.Run()
}
