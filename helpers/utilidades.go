package helpers

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"reflect"
)

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}

	client := &http.Client{}
	req, _ := http.NewRequest(trequest, url, b)

	defer func() {
		//Catch
		if r := recover(); r != nil {

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				beego.Error("Error reading response. ", err)
			}

			defer resp.Body.Close()
			mensaje, err := io.ReadAll(resp.Body)
			if err != nil {
				beego.Error("Error converting response. ", err)
			}
			bodyreq, err := io.ReadAll(req.Body)
			if err != nil {
				beego.Error("Error converting response. ", err)
			}
			respuesta := map[string]interface{}{"request": map[string]interface{}{"url": req.URL.String(), "header": req.Header, "body": bodyreq}, "body": mensaje, "statusCode": resp.StatusCode, "status": resp.Status}
			e, err := json.Marshal(respuesta)
			if err != nil {
				logs.Error(err)
			}
			json.Unmarshal(e, &target)
		}
	}()

	req.Header.Set("Authorization", "")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("accept", "*/*")

	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func DefaultTo[T any](value, defaultValue T) T {
	if reflect.ValueOf(value).IsZero() {
		return defaultValue
	} else {
		return value
	}
}

func DefaultToMapString(objMap map[string]any, key string, defaultValue any) any {
	if value, hasKey := objMap[key]; hasKey {
		return value
	} else {
		return defaultValue
	}
}
