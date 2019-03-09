package request

import (
	"encoding/json"
	"gitee.com/NotOnlyBooks/statistical_report/exception"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

type Validator interface {
	Validation() error
}

func parseParameter(c *gin.Context, tagName string) string {
	res := c.Param(tagName)

	if c.Query(tagName) != "" {
		res = c.Query(tagName)
	}

	return res
}

func BindQueryParams(c *gin.Context, a interface{}) exception.Exception {
	var (
		v reflect.Value
		t reflect.Type
	)

	if reflect.ValueOf(a).Kind() == reflect.Ptr {
		v = reflect.ValueOf(a).Elem()
		t = reflect.TypeOf(a).Elem()
	} else {
		v = reflect.ValueOf(a)
		t = reflect.TypeOf(a)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagName := field.Tag.Get("param")

		if tagName == "" {
			continue
		}

		param := parseParameter(c, tagName)

		if param == "" {
			continue
		}

		f := v.FieldByName(field.Name)

		switch field.Type.Kind() {
		case reflect.String:
			f.SetString(param)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uit, err := strconv.ParseUint(param, 10, 64); err != nil {
				return exception.ExceptionIllegalParameter
			} else {
				f.SetUint(uit)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if its, err := strconv.ParseInt(param, 10, 64); err != nil {
				return exception.ExceptionIllegalParameter
			} else {
				f.SetInt(its)
			}
		}
	}

	return nil
}

func Bind(c *gin.Context, a interface{}) exception.Exception {
	// first step bind url query parameter
	if excpt := BindQueryParams(c, a); excpt != nil {
		return excpt
	}

	// next step bind the body parameter
	buf, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		if err != io.EOF {
			zap.L().Warn("Request.Body.Bind.Error", zap.Error(err))
			return exception.ExceptionIllegalParameter
		}
	}

	if len(buf) > 0 {
		if err := json.Unmarshal(buf, a); err != nil {
			zap.L().Debug("Request.Body.Unmarshal.Error", zap.Error(err))
			return exception.ExceptionIllegalParameter
		}
	}

	// whether contain validation function
	if vFunc, ok := a.(Validator); !ok {
		panic("These is no implement validation func")
	} else {
		if err := vFunc.Validation(); err != nil {
			return exception.NewFromErr(exception.IllegalParameterCode, err)
		}
	}

	return nil
}
