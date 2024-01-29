package requestManager

import (
	"fmt"
	"strings"
)

type response struct {
	Success bool
	Status  string
	Message string
	Data    interface{}
}

// Formato de respuesta generalizado para entrega de respuesta de MID
//   - from: indica de que controlador va la info o de donde proviene el error
//   - method: POST, GET, PUT, DELETE
//   - success: si exitoso o no
//   - data: cuerpo de la respuesta
//
// Retorna:
//   - respuesta formateada
//   - status code
func MidResponseFormat(from string, method string, success bool, data interface{}) (response, int) {
	_method := strings.ToUpper(method)
	_status := 500
	_message := ""

	switch _method {
	case "POST":
		if success {
			_status = 201
			_message = "Registration successful"
		} else {
			_status = 400
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect data type or an invalid parameter", from)
		}
	case "GET":
		if success {
			_status = 200
			_message = "Request successful"
		} else {
			_status = 404
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect parameter or no record exist", from)
		}
	case "PUT":
		if success {
			_status = 200
			_message = "Update successful"
		} else {
			_status = 400
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect data type or an invalid parameter", from)
		}
	case "DELETE":
		if success {
			_status = 200
			_message = "Delete successful"
		} else {
			_status = 404
			_message = fmt.Sprintf("Error service %s: Request contains incorrect parameter", from)
		}
	}

	return response{
		Success: success,
		Status:  fmt.Sprintf("%d", _status),
		Message: _message,
		Data:    data,
	}, _status
}
