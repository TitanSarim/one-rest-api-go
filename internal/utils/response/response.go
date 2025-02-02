package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)


type Response struct {
    Status  string  `json:"status"`
    Error string `json:"error"`
}

const (
	StatusOk = "OK" 
	StatusError = "Error"
)

func WriteJson(res http.ResponseWriter, status int, data interface{}) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(data)
}

func GeneralError(err error) Response{
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var erMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			erMsgs = append(erMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			erMsgs = append(erMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(erMsgs, ", "),
	}
}