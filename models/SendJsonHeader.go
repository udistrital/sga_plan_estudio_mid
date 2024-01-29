package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
)

var global *context.Context

func GetHeader() (ctx *context.Context) {
	return global
}
func GetJson(urlp string, target interface{}) error {

	req, err := http.NewRequest("GET", urlp, nil)
	if err != nil {
		beego.Error("Error reading request. ", err)
	}

	//Se intenta acceder a cabecera, si no existe, se realiza peticion normal.

	defer func() {
		//Catch
		if r := recover(); r != nil {

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				beego.Error("Error reading response. ", err)
			}

			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(target)
		}
	}()

	//try
	header := GetHeader().Request.Header
	req.Header.Set("Authorization", header["Authorization"][0])
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Error reading response. ", err)
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
