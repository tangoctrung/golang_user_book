package helper

import "strings"

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObject struct{}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}

	return res
}

func BuildErrorsResponse(status bool, message string, errors string) Response {
	splittedError := strings.Split(errors, "\n")
	res := Response{
		Status:  status,
		Message: message,
		Errors:  splittedError,
		Data:    nil,
	}

	return res
}
