package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"],
        beego.ControllerComments{
            Method: "PostBaseStudyPlan",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"],
        beego.ControllerComments{
            Method: "GetStudyPlanVisualization",
            Router: "/:plan_id/estructura-visualizacion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"],
        beego.ControllerComments{
            Method: "GetPlanPorDependenciaVinculacionTercero",
            Router: "/dependencias-vinculacion-terceros/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_plan_estudio_mid/controllers:PlanEstudiosController"],
        beego.ControllerComments{
            Method: "PostGenerarDocumentoPlanEstudio",
            Router: "/generador-documentos-malla",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
