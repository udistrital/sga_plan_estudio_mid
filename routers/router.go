// @APIVersion 1.0.0
// @Title SGA MID - Plan de estudios
// @Description Microservicio MID del SGA que complementa los endpoints para el plan de estudios.
// @Contact
// @TermsOfServiceUrl
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/sga_plan_estudio_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/plan-estudios",
			beego.NSInclude(
				&controllers.PlanEstudiosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
