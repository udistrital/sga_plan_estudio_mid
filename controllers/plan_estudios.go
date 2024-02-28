package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_plan_estudio_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// PlanEstudiosController operations for Plan_estudios
type PlanEstudiosController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanEstudiosController) URLMapping() {
	c.Mapping("PostBaseStudyPlan", c.PostBaseStudyPlan)
	c.Mapping("GetStudyPlanVisualization", c.GetStudyPlanVisualization)
	c.Mapping("PostGenerarDocumentoPlanEstudio", c.PostGenerarDocumentoPlanEstudio)
	c.Mapping("GetPlanPorDependenciaVinculacionTercero", c.GetPlanPorDependenciaVinculacionTercero)
}

// PostBaseStudyPlan ...
// @Title PostBaseStudyPlan
// @Description create study plan
// @Param	body		body 	{}	true		"body for Plan_estudios content"
// @Success 201 {}
// @Failure 403 body is empty
// @router / [post]
func (c *PlanEstudiosController) PostBaseStudyPlan() {
	defer errorhandler.HandlePanic(&c.Controller)
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.PostBaseStudyPlan(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// GetStudyPlanVisualization ...
// @Title GetStudyPlanVisualization
// @Description get study plan data to the visualization
// @Param	id_plan		path	int	true	"Id del plan de estudio"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:plan_id/estructura-visualizacion [get]
func (c *PlanEstudiosController) GetStudyPlanVisualization() {
	defer errorhandler.HandlePanic(&c.Controller)
	idPlanString := c.Ctx.Input.Param(":plan_id")
	idPlan, errId := strconv.ParseInt(idPlanString, 10, 64)
	if errId != nil || idPlan <= 0 {
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, errId.Error())
		c.ServeJSON()
	} else {
		resultado := services.GetStudyPlanVisualization(idPlan)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
		c.ServeJSON()
	}
}

// PostGenerarDocumentoPlanEstudio ...
// @Title PostGenerarDocumentoPlanEstudio
// @Description Genera un documento PDF del plan de estudio
// @Param	body		body 	{}	true		"body Datos del plan de estudio content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /generador-documentos-malla [post]
func (c *PlanEstudiosController) PostGenerarDocumentoPlanEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.PostGenerarDocumentoPlanEstudio(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// GetPlanPorDependenciaVinculacionTercero ...
// @Title GetPlanPorDependenciaVinculacionTercero
// @Description get plan de estudio por DependenciaId de vinculación de tercero, verificando cargo
// @Param	body		body 	{}	true		"body Datos del plan de estudio content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /dependencias-vinculacion-terceros/:tercero_id [get]
func (c *PlanEstudiosController) GetPlanPorDependenciaVinculacionTercero() {
	defer errorhandler.HandlePanic(&c.Controller)
	terceroIdStr := c.Ctx.Input.Param(":tercero_id")
	terceroId, errId := strconv.ParseInt(terceroIdStr, 10, 64)
	if errId != nil || terceroId <= 0 {
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, errId.Error())
		c.ServeJSON()
	} else {
		resultado := services.GetPlanPorDependenciaVinculacionTercero(terceroId)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
		c.ServeJSON()
	}

}
