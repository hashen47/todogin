package handlers

import (
	"fmt"
	"errors"
	"reflect"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Status string 
const (
	OK   Status = "success"
	FAIL Status = "fail"
)
type ErrsMap map[string]map[string]string

func NewResp(status Status, data map[string]any, err error, errs ErrsMap) gin.H {
	resp := gin.H{}
	resp["status"] = string(status)
	resp["data"]   = data 

	if err != nil {
		resp["error"] = err.Error()
	}

	if len(errs) > 0 {
		resp["errors"] = errs
	}

	return resp 
}

// don't use nested structures here, (mean don't use with 'dive' validation option)
func GetErrorMsgs(obj any, err error) (ErrsMap, error) {
	errs := make(ErrsMap, 0)

	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return errs, fmt.Errorf("Invalid json type for field: %q, expected type: %q", unmarshalTypeErr.Field, unmarshalTypeErr.Type)
	}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			jsonName := ""
			field, ok := reflect.TypeOf(obj).FieldByName(e.Field())
			if ok {
				jsonName = field.Tag.Get("json")
			}

			errs[jsonName] = make(map[string]string, 0)

			switch e.Tag() {
			case "required":
				errs[jsonName][e.Tag()] = fmt.Sprintf("%s is required", jsonName)
			case "min":
				errs[jsonName][e.Tag()] = fmt.Sprintf("%s should at least %s (current: %s)", jsonName, e.Param(), e.Value())
			case "max":
				errs[jsonName][e.Tag()] = fmt.Sprintf("%s cannot be exceed %s (current: %s)", jsonName, e.Param(), e.Value())
			case "email":
				errs[jsonName][e.Tag()] = fmt.Sprintf("invalid email value for %s (value: %s)", jsonName, e.Value())
			default:
				panic("(GetErrorMsgs) not implemented yet")
			}
		}
	}

	return errs, nil
}
