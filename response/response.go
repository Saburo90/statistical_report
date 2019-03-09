package response

import (
	"gitee.com/NotOnlyBooks/statistical_report/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func getDataParameter(data []interface{}) interface{} {
	if len(data) == 0 {
		return nil
	}

	return data[0]
}

func SuccessResp(contexts *gin.Context, data ...interface{}) {
	contexts.JSON(http.StatusOK, Response{
		Code: exception.SuccessCode,
		Msg:  exception.Success.GetMsg(),
		Data: getDataParameter(data),
	})
}

func ThrowException(contexts *gin.Context, e exception.Exception, data ...interface{}) {
	status := http.StatusOK
	if e == exception.ExceptionIllegalSign {
		status = http.StatusUnauthorized
	}

	contexts.AbortWithStatusJSON(status, Response{
		Code: e.GetCode(),
		Msg:  e.GetMsg(),
		Data: getDataParameter(data),
	})
}
