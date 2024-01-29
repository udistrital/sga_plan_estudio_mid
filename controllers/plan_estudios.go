package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	request "github.com/udistrital/sga_mid_plan_estudios/models"
	"github.com/udistrital/sga_mid_plan_estudios/process/plan_estudio_visualizacion_documento"
	"github.com/udistrital/sga_mid_plan_estudios/utils"
	requestmanager "github.com/udistrital/sga_mid_plan_estudios/utils/requestManager"
	"reflect"
	"sga_mid_plan_estudios/helpers"
	"sort"
	"strconv"
	"strings"
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
	var studyPlanRequest map[string]interface{}
	const editionApprovalStatus = 1

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyPlanRequest); err == nil {
		if status, errStatus := getApprovalStatus(editionApprovalStatus); errStatus == nil {
			studyPlanRequest["EstadoAprobacionId"] = status.(map[string]interface{})
			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "EspaciosSemestreDistribucion"); errMap == nil {
				if newString != "" {
					studyPlanRequest["EspaciosSemestreDistribucion"] = newString
				}
			}

			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "ResumenPlanEstudios"); errMap == nil {
				if newString != "" {
					studyPlanRequest["ResumenPlanEstudios"] = newString
				}
			}

			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "SoporteDocumental"); errMap == nil {
				if newString != "" {
					studyPlanRequest["SoporteDocumental"] = newString
				}
			}

			if newPlan, errPlan := createStudyPlan(studyPlanRequest); errPlan == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = map[string]interface{}{
					"Success": true, "Status": "201",
					"Message": "Created",
					"Data":    newPlan,
				}
			} else {
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{
					"Success": false, "Status": "400",
					"Message": "Error al crear el plan de estudios",
				}
			}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{
				"Success": false, "Status": "404",
				"Message": "Estado aprobación del plan de estudios no encontrado",
			}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"CreacionPlanEstudioBase", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
	c.ServeJSON()
}

func getApprovalStatus(id int) (any, error) {
	var resStudyPlan interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		"estado_aprobacion?" + "query=id:" + fmt.Sprintf("%v", id)
	if errPlan := request.GetJson(urlStudyPlan, &resStudyPlan); errPlan == nil {
		if resStudyPlan.(map[string]interface{})["Data"] != nil {
			status := resStudyPlan.(map[string]interface{})["Data"].([]any)[0]
			return status, nil
		} else {
			return nil, fmt.Errorf("PlanEstudiosService No se encuentra el estado aprobación requerido")
		}
	} else {
		return nil, errPlan
	}
}

func createStudyPlan(studyPlanBody map[string]interface{}) (map[string]interface{}, error) {
	var newStudyPlan map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		"plan_estudio"
	if errNewPlan := helpers.SendJson(urlStudyPlan, "POST", &newStudyPlan, studyPlanBody); errNewPlan == nil {
		return newStudyPlan["Data"].(map[string]interface{}), nil
	} else {
		return newStudyPlan, fmt.Errorf("PlanEstudiosService Error creando plan de estudios")
	}
}

func map2StringFieldStudyPlan(body map[string]any, fieldName string) (string, error) {
	if reflect.TypeOf(body[fieldName]).Kind() == reflect.Map {
		if stringNew, errMS := utils.Map2String(body[fieldName].(map[string]interface{})); errMS == nil {
			return stringNew, nil
		} else {
			return "", errMS
		}
	} else {
		return "", nil
	}
}

// GetStudyPlanVisualization ...
// @Title GetStudyPlanVisualization
// @Description get study plan data to the visualization
// @Param	id_plan		path	int	true	"Id del plan de estudio"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:plan_id/estructura-visualizacion [get]
func (c *PlanEstudiosController) GetStudyPlanVisualization() {

	failureAsn := map[string]interface{}{
		"Success": false,
		"Status":  "404",
		"Message": "Error service GetStudyPlanVisualization: The request contains an incorrect parameter or no record exist",
		"Data":    nil}
	successAns := map[string]interface{}{
		"Success": true,
		"Status":  "200",
		"Message": "Query successful",
		"Data":    nil}
	idPlanString := c.Ctx.Input.Param(":plan_id")
	idPlan, errId := strconv.ParseInt(idPlanString, 10, 64)
	if errId != nil || idPlan <= 0 {
		if errId == nil {
			errId = fmt.Errorf("id_plan: %d <= 0", idPlan)
		}
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errId.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	var resStudyPlan map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		fmt.Sprintf("plan_estudio/%v", idPlan)
	errPlan := request.GetJson(urlStudyPlan, &resStudyPlan)

	if errPlan == nil && resStudyPlan["Success"] == true {
		planData := resStudyPlan["Data"].(map[string]interface{})

		classificationsData, errorClass := getClassificationData()
		if errorClass != nil {
			logs.Error(errorClass.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = errorClass.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}

		if planData["EsPlanEstudioPadre"] == true {
			visualizationData, errPlanVisualization := getParentStudyPlanVisualization(planData, classificationsData)
			if errPlanVisualization == nil {
				c.Ctx.Output.SetStatus(200)
				successAns["Data"] = visualizationData
				c.Data["json"] = successAns
				c.ServeJSON()
			} else {
				logs.Error(errPlanVisualization.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = errPlanVisualization.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		} else {
			visualizationData, errPlanVisualization := getChildStudyPlanVisualization(planData, classificationsData)
			if errPlanVisualization == nil {
				c.Ctx.Output.SetStatus(200)
				successAns["Data"] = visualizationData
				c.Data["json"] = successAns
				c.ServeJSON()
			} else {
				logs.Error(errPlanVisualization.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = errPlanVisualization.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		}
	} else {
		if errPlan == nil {
			errPlan = fmt.Errorf("PlanEstudioService: %v", resStudyPlan["Message"])
		}
		logs.Error(errPlan.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errPlan.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
}

func getChildStudyPlanVisualization(studyPlanData map[string]interface{}, classificationsData []interface{}) (map[string]interface{}, error) {
	var facultyName string
	var totalPlanData []map[string]interface{}
	academies := map[string]map[string]interface{}{}

	if studyPlanData["ProyectoAcademicoId"] != nil {
		projectCData, projectErr := utils.GetProyectoCurricular(int(studyPlanData["ProyectoAcademicoId"].(float64)))
		if projectErr == nil {
			facultyData, errFacultad := utils.GetFacultadDelProyectoC(fmt.Sprintf("%v", projectCData["id_oikos"]))
			if errFacultad == nil {
				facultyName = fmt.Sprintf("%v", facultyData["Nombre"])
			}

			planData, errorPlan := getPlanVisualization(studyPlanData, 1,
				fmt.Sprintf("%v", projectCData["id_snies"]), classificationsData, academies)
			if errorPlan == nil {
				totalPlanData = append(totalPlanData, planData)
			} else {
				logs.Error(errorPlan.Error())
				return nil, errorPlan
			}
		} else {
			return nil, projectErr
		}
	} else {
		return nil, fmt.Errorf("without ProyectoAcademicoId")
	}
	dataResult := map[string]any{
		"Nombre":   studyPlanData["Nombre"],
		"Facultad": helpers.DefaultTo(facultyName, ""),
		"Planes":   totalPlanData,
		"Escuelas": academies,
	}
	return dataResult, nil
}

func getPlanVisualization(studyPlanData map[string]interface{}, orderNumb int, snies string, classificationsData []interface{}, academies map[string]map[string]interface{}) (map[string]interface{}, error) {
	var resolution string
	var semesterDistribution map[string]interface{}
	periodInfoTotal := []map[string]interface{}{}
	summary := map[string]interface{}{}

	if studyPlanData["NumeroResolucion"] != nil && studyPlanData["AnoResolucion"] != nil {
		resolution = fmt.Sprintf("%v de %v", studyPlanData["NumeroResolucion"], studyPlanData["AnoResolucion"])
	} else {
		resolution = ""
	}

	if semesterDistributionData, semesterDataOk := studyPlanData["EspaciosSemestreDistribucion"]; semesterDataOk && semesterDistributionData != nil {
		if reflect.TypeOf(semesterDistributionData).Kind() == reflect.String {
			if err := json.Unmarshal([]byte(semesterDistributionData.(string)), &semesterDistribution); err == nil {
				spaceVisualizationData, summaryData, errorSpaceVisualization := semesterDistribution2SpacesVisualization(
					semesterDistribution, classificationsData, academies)
				if errorSpaceVisualization == nil {
					periodInfoTotal = spaceVisualizationData
					summary = summaryData
				} else {
					return nil, errorSpaceVisualization
				}
			}
		}
	}

	planData := map[string]interface{}{
		"Orden":        orderNumb,
		"Nombre":       helpers.DefaultTo(studyPlanData["Nombre"], ""),
		"Resolucion":   resolution,
		"Creditos":     studyPlanData["TotalCreditos"],
		"Snies":        helpers.DefaultTo(snies, ""),
		"PlanEstudio":  helpers.DefaultTo(studyPlanData["Codigo"], ""),
		"InfoPeriodos": periodInfoTotal,
		"Resumen":      summary,
	}
	return planData, nil
}

func semesterDistribution2SpacesVisualization(spaceSemesterDistribution map[string]interface{}, classificationsData []interface{}, academies map[string]map[string]interface{}) ([]map[string]any, map[string]any, error) {
	var periodOrder = 1
	var totalSpaceVisualizationData []map[string]interface{}
	var totalPeriodData []map[string]interface{}
	summary := map[string]interface{}{}

	totalOB := 0
	totalOC := 0
	totalEI := 0
	totalEE := 0

	// Sort spaceSemesterDistribution
	keysSpace := make([]string, 0, len(spaceSemesterDistribution))
	for k := range spaceSemesterDistribution {
		keysSpace = append(keysSpace, k)
	}
	sort.Strings(keysSpace)

	// Iterate every semester
	for _, kSpace := range keysSpace {
		semesterV := spaceSemesterDistribution[kSpace]
		totalCredits := 0
		totalSpaceVisualizationData = []map[string]interface{}{}
		if spaces, spaceOk := semesterV.(map[string]interface{})["espacios_academicos"]; spaceOk && spaces != nil {
			if reflect.TypeOf(spaces).Kind() == reflect.Array || reflect.TypeOf(spaces).Kind() == reflect.Slice {
				// Iterate every space
				for _, spaceV := range spaces.([]interface{}) {
					if reflect.TypeOf(spaceV).Kind() == reflect.Map {
						//	Get space data
						spaceData := utils.MapValues(spaceV.(map[string]interface{}))
						var spaceId string
						for _, spaceField := range spaceData {
							if spaceField.(map[string]interface{})["Id"] != nil {
								spaceId = fmt.Sprintf("%v",
									spaceField.(map[string]interface{})["Id"])
								spaceVisualizationData, spaceVisualizationErr := getSpaceVisualizationData(
									spaceId, classificationsData, academies)

								if spaceVisualizationErr != nil {
									return nil, nil, spaceVisualizationErr
								} else {
									totalCredits = totalCredits + int(spaceVisualizationData["Creditos"].(float64))
									if spaceVisualizationData["Clasificacion"] == "OB" {
										totalOB = totalOB + int(spaceVisualizationData["Creditos"].(float64))
									} else if spaceVisualizationData["Clasificacion"] == "OC" {
										totalOC = totalOC + int(spaceVisualizationData["Creditos"].(float64))
									} else if spaceVisualizationData["Clasificacion"] == "EI" {
										totalEI = totalEI + int(spaceVisualizationData["Creditos"].(float64))
									} else if spaceVisualizationData["Clasificacion"] == "EE" {
										totalEE = totalEE + int(spaceVisualizationData["Creditos"].(float64))
									}

									totalSpaceVisualizationData = append(
										totalSpaceVisualizationData,
										spaceVisualizationData)
								}
							}
						}
					}
				}
			}
		}
		periodData := map[string]interface{}{
			"Orden":         periodOrder,
			"Espacios":      totalSpaceVisualizationData,
			"TotalCreditos": totalCredits,
		}

		totalPeriodData = append(totalPeriodData, periodData)
		periodOrder++
	}
	summary = map[string]interface{}{
		"OB": totalOB,
		"OC": totalOC,
		"EI": totalEI,
		"EE": totalEE,
	}
	return totalPeriodData, summary, nil
}

func getSpaceVisualizationData(academicSpaceId string, classificationsData []interface{}, academies map[string]map[string]interface{}) (map[string]interface{}, error) {
	var academicSpace map[string]interface{}
	url := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
		fmt.Sprintf("espacio-academico/%v", academicSpaceId)

	academicSpaceError := request.GetJson(url, &academicSpace)
	if academicSpaceError != nil || academicSpace == nil || academicSpace["Success"] == false {
		return nil, fmt.Errorf("EspaciosAcademicosService: %v", academicSpace["Message"])
	}

	academicSpaceData := academicSpace["Data"].(map[string]interface{})
	var hoursDistributionData map[string]interface{}
	if hoursDistribution, hoursOk := academicSpaceData["distribucion_horas"]; hoursOk {
		hoursDistributionData = hoursDistribution.(map[string]interface{})
	} else {
		hoursDistributionData = map[string]interface{}{
			"HTA": 0,
			"HTC": 0,
			"HTD": 0}
	}
	classificationCode, classificationErr := getClassificationVisualizationData(
		academicSpaceData["clasificacion_espacio_id"].(float64),
		classificationsData)
	if classificationErr != nil {
		classificationCode = map[string]interface{}{
			"Nombre":            "",
			"CodigoAbreviacion": ""}
	}

	// Prerequisites
	prerequisitesCode := []string{}
	if reflect.TypeOf(academicSpaceData["espacios_requeridos"]).Kind() == reflect.Array || reflect.TypeOf(academicSpaceData["espacios_requeridos"]).Kind() == reflect.Slice {
		for _, prerequisiteId := range academicSpaceData["espacios_requeridos"].([]interface{}) {
			var prerequisiteResponse map[string]interface{}
			url := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
				fmt.Sprintf("espacio-academico/%v", prerequisiteId)
			prerequisiteError := request.GetJson(url, &prerequisiteResponse)

			if prerequisiteError != nil || prerequisiteResponse["Success"] == false {
				return nil, fmt.Errorf("EspaciosAcademicosService: Prerequisite not found. %v",
					academicSpace["Message"])
			}
			prerequisiteData := prerequisiteResponse["Data"].(map[string]interface{})
			prerequisitesCode = append(prerequisitesCode,
				fmt.Sprintf("%v", prerequisiteData["codigo"]))
		}
	}

	// Get academic space colors
	var groupingSpace map[string]interface{}
	var groupingSpaceData map[string]interface{}
	groupingSpaceId, groupingSpaceOk := academicSpaceData["agrupacion_espacios_id"]
	if groupingSpaceOk {
		// Get academic
		url = "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
			fmt.Sprintf("agrupacion-espacios/%v", groupingSpaceId)

		groupingSpaceError := request.GetJson(url, &groupingSpace)
		if groupingSpaceError != nil || groupingSpace["Success"] == false {
			groupingSpaceData = map[string]interface{}{
				"color_hex": "#FFFFFF",
				"nombre":    "---",
			}
		} else {
			groupingSpaceData = groupingSpace["Data"].(map[string]interface{})
		}

		// Set academies
		specificAcademic, specificAcademicOk := academies[fmt.Sprintf("%v", groupingSpaceId)]
		if specificAcademicOk {
			academies[fmt.Sprintf("%v", groupingSpaceId)]["Cantidad"] = specificAcademic["Cantidad"].(int) + 1
		} else {
			academies[fmt.Sprintf("%v", groupingSpaceId)] = map[string]any{
				"Cantidad": 1,
				"Color":    groupingSpaceData["color_hex"],
				"Nombre":   groupingSpaceData["nombre"],
			}
		}
	} else {
		groupingSpaceId = "sinEscuela"

		groupingSpaceData = map[string]interface{}{
			"color_hex": "#FFFFFF",
			"nombre":    "---",
		}

		// Set academies
		specificAcademic, specificAcademicOk := academies[fmt.Sprintf("%v", groupingSpaceId)]
		if specificAcademicOk {
			academies[fmt.Sprintf("%v", groupingSpaceId)]["Cantidad"] = specificAcademic["Cantidad"].(int) + 1
		} else {
			academies[fmt.Sprintf("%v", groupingSpaceId)] = map[string]any{
				"Cantidad": 1,
				"Color":    groupingSpaceData["color_hex"],
				"Nombre":   groupingSpaceData["nombre"],
			}
		}
	}

	fontcolor := getFontColor(fmt.Sprintf("%v", groupingSpaceData["color_hex"]))
	academies[fmt.Sprintf("%v", groupingSpaceId)]["TxtColor"] = fontcolor

	spaceResult := map[string]interface{}{
		"Codigo":             academicSpaceData["codigo"],
		"Nombre":             academicSpaceData["nombre"],
		"Creditos":           academicSpaceData["creditos"],
		"Prerequisitos":      prerequisitesCode,
		"HTD":                hoursDistributionData["HTD"],
		"HTC":                hoursDistributionData["HTC"],
		"HTA":                hoursDistributionData["HTA"],
		"Clasificacion":      classificationCode["CodigoAbreviacion"],
		"Escuela":            groupingSpaceId,
		"EscuelaColor":       groupingSpaceData["color_hex"],
		"EscuelaColorFuente": fontcolor,
	}

	return spaceResult, nil
}

func getClassificationVisualizationData(idClassification float64, classifications []interface{}) (map[string]interface{}, error) {
	// Get class by idClassification
	for _, classData := range classifications {
		if idClassification == classData.(map[string]interface{})["Id"].(float64) {
			result := map[string]interface{}{
				"Nombre":            classData.(map[string]interface{})["Nombre"],
				"CodigoAbreviacion": classData.(map[string]interface{})["CodigoAbreviacion"]}
			return result, nil
		}
	}
	return nil, fmt.Errorf("classification not found")
}

func getClassificationData() ([]interface{}, error) {
	classId := 51
	var spaceClassResult map[string]interface{}

	spaceClassErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,Nombre,CodigoAbreviacion", classId), &spaceClassResult)
	if spaceClassErr != nil || fmt.Sprintf("%v", spaceClassResult) == "[map[]]" {
		if spaceClassErr == nil {
			spaceClassErr = fmt.Errorf("ParametroService: query for clases is empty")
		}
		logs.Error(spaceClassErr.Error())
		return nil, spaceClassErr
	}

	if classificationsData, classOk := spaceClassResult["Data"]; classOk {
		return classificationsData.([]interface{}), nil
	} else {
		return nil, fmt.Errorf("ParametroService: Without data to space classifications")
	}
}

func getParentStudyPlanVisualization(studyPlanData map[string]interface{}, classificationsData []interface{}) (map[string]interface{}, error) {
	var facultyName string
	var totalPlanData []map[string]interface{}
	var orderPlan map[string]interface{}
	academies := map[string]map[string]interface{}{}

	if studyPlanData["ProyectoAcademicoId"] != nil {
		projectCData, projectErr := utils.GetProyectoCurricular(int(studyPlanData["ProyectoAcademicoId"].(float64)))
		if projectErr == nil {
			facultyData, errFacultad := utils.GetFacultadDelProyectoC(fmt.Sprintf("%v", projectCData["id_oikos"]))
			if errFacultad == nil {
				facultyName = fmt.Sprintf("%v", facultyData["Nombre"])
			}

			planProjectData, errorPlanProject := getPlanProjectByParent(studyPlanData["Id"].(float64))
			if errorPlanProject == nil {
				if orderPlanData, orderDataOk := planProjectData["OrdenPlan"]; orderDataOk && orderPlanData != nil {
					if reflect.TypeOf(orderPlanData).Kind() == reflect.String {
						if err := json.Unmarshal([]byte(orderPlanData.(string)), &orderPlan); err == nil {
							// Sort plans
							keysPlans := make([]string, 0, len(orderPlan))
							for k := range orderPlan {
								keysPlans = append(keysPlans, k)
							}
							sort.SliceStable(keysPlans, func(i, j int) bool {
								return orderPlan[keysPlans[i]].(map[string]any)["Orden"].(float64) < orderPlan[keysPlans[j]].(map[string]any)["Orden"].(float64)
							})

							for _, kPlan := range keysPlans {
								planV := orderPlan[kPlan]
								if idChildPlan, childPlanError := planV.(map[string]interface{})["Id"]; childPlanError {
									var resChildStudyPlan map[string]interface{}
									urlChildStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
										fmt.Sprintf("plan_estudio/%v", idChildPlan)

									errChildPlan := request.GetJson(urlChildStudyPlan, &resChildStudyPlan)

									if errChildPlan == nil && resChildStudyPlan["Success"] == true {
										childStudyPlanData := resChildStudyPlan["Data"].(map[string]interface{})

										projectCData, projectErr := utils.GetProyectoCurricular(
											int(childStudyPlanData["ProyectoAcademicoId"].(float64)))

										if projectErr != nil {
											return nil, projectErr
										}

										childPlanData, errorPlan := getPlanVisualization(
											childStudyPlanData,
											int(planV.(map[string]interface{})["Orden"].(float64)),
											fmt.Sprintf("%v", projectCData["id_snies"]),
											classificationsData,
											academies)
										if errorPlan == nil {
											totalPlanData = append(totalPlanData, childPlanData)
										} else {
											return nil, errorPlan
										}
									} else {
										if errChildPlan == nil {
											errChildPlan = fmt.Errorf("PlanEstudioService: %v", resChildStudyPlan["Message"])
										}
										logs.Error(errChildPlan.Error())
										return nil, errChildPlan
									}
								} else {
									return nil, fmt.Errorf("error getting id child plan")
								}
							}
						} else {
							return nil, fmt.Errorf("error getting plan order, OrdenPlan field")
						}
					}
				}
			} else {
				logs.Error(errorPlanProject.Error())
				return nil, errorPlanProject
			}
		} else {
			return nil, projectErr
		}
	} else {
		return nil, fmt.Errorf("without ProyectoAcademicoId")
	}
	dataResult := map[string]any{
		"Nombre":   studyPlanData["Nombre"],
		"Facultad": helpers.DefaultTo(facultyName, ""),
		"Planes":   totalPlanData,
		"Escuelas": academies,
	}
	return dataResult, nil
}

func getPlanProjectByParent(parentId float64) (map[string]any, error) {
	var resStudyPlanProject map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		fmt.Sprintf("plan_estudio_proyecto_academico?query=activo:true,PlanEstudioId:%v", parentId)
	errPlan := request.GetJson(urlStudyPlan, &resStudyPlanProject)

	if errPlan == nil && resStudyPlanProject["Success"] == true && resStudyPlanProject["Status"] == "200" {
		studyPlanProjectData := resStudyPlanProject["Data"].([]interface{})

		if len(studyPlanProjectData) > 0 {
			return studyPlanProjectData[0].(map[string]interface{}), nil
		} else {
			return nil, fmt.Errorf("PlanEstudioService: Without data in plan_estudio_proyecto_academico")
		}
	} else {
		return nil, fmt.Errorf("PlanEstudioService: Error in request plan_estudio_proyecto_academico")
	}
}

func getFontColor(backgroundColor string) string {
	hexString := strings.Replace(backgroundColor, "#", "", 1)

	var rgb [3]int
	values, err := strconv.ParseInt(hexString, 16, 32)

	if err != nil {
		return "#838383"
	}

	rgb[0] = int(values >> 16)
	rgb[1] = int((values >> 8) & 0xFF)
	rgb[2] = int(values & 0xFF)

	brightness := (float64(rgb[0])*0.299 + float64(rgb[1])*587 + float64(rgb[2])*114) / 255

	if brightness > 0.5 {
		return "#000000"
	} else {
		return "#FFFFFF"
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
	var data map[string]interface{}

	if parseErr := json.Unmarshal(c.Ctx.Input.RequestBody, &data); parseErr == nil {
		plans := data["Planes"].([]any)
		if len(plans) > 0 {
			// sort plans by order field
			sort.Slice(plans, func(i, j int) bool {
				return plans[i].(map[string]any)["Orden"].(float64) < plans[j].(map[string]any)["Orden"].(float64)
			})

			//	sort periods of each plan by the order field
			for _, plan := range plans {
				infoPeriods := plan.(map[string]any)["InfoPeriodos"].([]any)
				sort.Slice(infoPeriods, func(i, j int) bool {
					return infoPeriods[i].(map[string]any)["Orden"].(float64) < infoPeriods[j].(map[string]any)["Orden"].(float64)
				})
			}

			// sort by academy
			academies := data["Escuelas"].(map[string]any)
			for _, plan := range plans {
				infoPeriods := plan.(map[string]any)["InfoPeriodos"].([]any)
				for _, period := range infoPeriods {
					spaces := period.(map[string]any)["Espacios"].([]any)
					sort.Slice(spaces, func(i, j int) bool {
						academyA := academies[fmt.Sprintf("%v", spaces[i].(map[string]any)["Escuela"].(string))]
						academyB := academies[fmt.Sprintf("%v", spaces[j].(map[string]any)["Escuela"].(string))]
						return academyA.(map[string]any)["Cantidad"].(float64) > academyB.(map[string]any)["Cantidad"].(float64)
					})
				}
			}
		}

		pdf := plan_estudio_visualizacion_documento.GenerateStudyPlanDocument(data)

		if pdf.Err() {
			logs.Error("Failed creating PDF report: %s\n", pdf.Error())
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{
				"Success": false, "Status": "400",
				"Message": "Error al generar el documento del plan de estudios",
			}
		}

		if pdf.Ok() {
			encodedFile := utils.EncodePDF(pdf)
			c.Data["json"] = map[string]interface{}{
				"Success": true,
				"Status":  "200",
				"Message": "Query successful",
				"Data":    encodedFile}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"PostGenerarDocumentoPlanEstudio", "POST", false, parseErr.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
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
	var plans []map[string]interface{}
	terceroIdStr := c.Ctx.Input.Param(":tercero_id")
	terceroId, errId := strconv.ParseInt(terceroIdStr, 10, 64)
	if errId != nil || terceroId <= 0 {
		if errId == nil {
			errId = fmt.Errorf("tercero_id: %d <= 0", terceroId)
		}
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"Success": false, "Status": "404",
			"Message": "Error service GetPlanPorDependenciaVinculacionTercero: The request contains an incorrect parameter",
			"Data":    errId.Error()}
		c.ServeJSON()
		return
	}

	// 1. Consultar dependencias por tercero
	/*
		consulta vinculación tercero y check resultado válido
		DependenciaId__gt:0 -> que tenga id mayor que cero
		CargoId__in:312|320 -> parametrosId: 312: JEFE OFICINA, 320: Asistente Dependencia
	*/
	var estadoVinculacion []map[string]interface{}
	estadoVinculacionErr := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+
		fmt.Sprintf("vinculacion?query=Activo:true,DependenciaId__gt:0,CargoId__in:312|320,tercero_principal_id:%v", terceroIdStr), &estadoVinculacion)
	if estadoVinculacionErr != nil || fmt.Sprintf("%v", estadoVinculacion) == "[map[]]" {
		if estadoVinculacionErr == nil {
			estadoVinculacionErr = fmt.Errorf("vinculacion is empty: %v", estadoVinculacion)
		}
		logs.Error(estadoVinculacionErr.Error())
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"Success": false, "Status": "404",
			"Message": "Error service GetPlanPorDependenciaVinculacionTercero: Tercero relationship is empty",
			"Data":    errId.Error()}
		c.ServeJSON()
		return
	}
	/*
		preparar lista de dependencias, normalmente será una, pero se espera soportar varias por tercero
	*/
	var dependencias []int64
	for _, vinculacion := range estadoVinculacion {
		dependencias = append(dependencias, int64(vinculacion["DependenciaId"].(float64)))
	}

	if len(dependencias) == 0 {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"Success": false, "Status": "404",
			"Message": "Error service GetPlanPorDependenciaVinculacionTercero: Tercero without dependencies",
			"Data":    fmt.Errorf("tercero without dependencies: %v", dependencias).Error()}
		c.ServeJSON()
		return
	}
	// 2. Consultar planes, revisando cuales despendencias si direron resultado
	for _, dependencia := range dependencias {
		var resStudyPlan map[string]interface{}
		urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
			fmt.Sprintf("plan_estudio?query=Activo:true,ProyectoAcademicoId:%v", dependencia)
		errPlan := request.GetJson(urlStudyPlan, &resStudyPlan)

		if errPlan == nil && resStudyPlan["Success"] == true {
			planData := resStudyPlan["Data"]
			if len(planData.([]interface{})) > 0 {
				// 3. Validar que las dependecias esten en proyecto academico, si no están se acepta,
				// si están se valida que la dependencia sea del tipo proyecto curricular
				// TipoDependenciaId: 1 -> PROYECTO CURRICULAR
				var resProject []map[string]interface{}
				urlStudyPlan := "http://" + beego.AppConfig.String("OikosService") +
					fmt.Sprintf("dependencia_tipo_dependencia?query=DependenciaId:%v&fields=TipoDependenciaId", dependencia)
				errProject := request.GetJson(urlStudyPlan, &resProject)
				if errProject == nil && resProject != nil && fmt.Sprintf("%v", resProject) != "[map[]]" {
					typeDependencieId, _ := strconv.ParseInt(fmt.Sprintf("%v", resProject[0]["TipoDependenciaId"]), 10, 64)
					if typeDependencieId == 1 {
						plans = append(plans, planData.([]map[string]interface{})...)
					}
				}
			}
		}
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  "200",
		"Message": "Query successful",
		"Data":    plans}
	c.ServeJSON()
}
