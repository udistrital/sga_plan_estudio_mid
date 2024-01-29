// Dependencias query functions.
// Funciones generalizadas para consultar los servicios de
// homologación dependencias o proyeto curricular y obtener los
// regisros resultantes

package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetFacultadDelProyectoC(projectIdOikos string) (map[string]any, error) {
	var facultadResponse map[string]interface{}

	facultadErr := request.GetJson(
		"http://"+beego.AppConfig.String("OikosService")+
			fmt.Sprintf("dependencia/get_dependencias_padres_by_id/%v", projectIdOikos),
		&facultadResponse)
	if facultadErr == nil && facultadResponse["Type"] == "success" {
		dependencias := facultadResponse["Body"].([]interface{})
		if len(dependencias) > 0 {
			var dependenciaPadre int
			for i := len(dependencias) - 1; i >= 0; i-- {
				if fmt.Sprintf("%v", dependencias[i].(map[string]interface{})["Id"]) == projectIdOikos {
					dependenciaPadre = int(dependencias[i].(map[string]interface{})["Padre"].(float64))
				}

				if dependenciaPadre > 0 && int(dependencias[i].(map[string]interface{})["Id"].(float64)) == dependenciaPadre {
					return dependencias[i].(map[string]interface{}), nil
				}
			}
		}
		return nil, fmt.Errorf("Facultad no encontrada")
	} else {
		return nil, fmt.Errorf("Facultad no encontrada")
	}
}

func GetProyectoCurricular(proyectoId int) (map[string]any, error) {
	var proyectoResponse map[string]interface{}

	proyectoErr := request.GetJsonWSO2(
		"http://"+beego.AppConfig.String("HomologacionDependenciaService")+
			fmt.Sprintf("proyecto_curricular_cod_proyecto/%v", proyectoId),
		&proyectoResponse)
	if proyectoErr == nil && fmt.Sprintf("%v", proyectoResponse) != "map[homologacion:map[]]" && fmt.Sprintf("%v", proyectoResponse) != "map[]]" {
		homologacionData := proyectoResponse["homologacion"].(map[string]interface{})
		proyectoData := map[string]interface{}{
			"proyecto_curricular_nombre": fmt.Sprintf("%v", homologacionData["proyecto_snies"]),
			"id_oikos":                   homologacionData["id_oikos"],
			"id_snies":                   homologacionData["id_snies"],
			"id_argo":                    homologacionData["id_argo"],
		}
		return proyectoData, nil
	} else {
		return nil, fmt.Errorf("Proyecto curricular homologación no encontrado")
	}
}
