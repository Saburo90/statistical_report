package mini_statis

import (
	"gitee.com/NotOnlyBooks/statistical_report/exception"
	"gitee.com/NotOnlyBooks/statistical_report/handler/mini_statis/interior"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"gitee.com/NotOnlyBooks/statistical_report/request"
	"gitee.com/NotOnlyBooks/statistical_report/response"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type MiniReq struct {
	Operator    string `json:"operator"`
	OperateTime int64  `json:"operateTime"`
	ClientIP    string `json:"clientIP"`
	Sign        string `json:"sign"`
}

func (req *MiniReq) Validation() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Operator, validation.Required, validation.In("admin", "saburo", "system")),
		validation.Field(&req.OperateTime, validation.Required),
		validation.Field(&req.ClientIP, validation.Required, validation.Match(regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"))),
		validation.Field(&req.Sign, validation.Required),
	)
}

func GetRoamMiniWxAPIDataHandler(c *gin.Context) {
	req := &MiniReq{}

	if excpt := request.Bind(c, req); excpt != nil {
		response.ThrowException(c, exception.ExceptionIllegalParameter)
		return
	}

	if req.Operator == "" || req.OperateTime == 0 || req.Sign == "" {
		response.ThrowException(c, exception.ExceptionIllegalParameter)
		return
	}

	// out of 5 minute throw exception
	if req.OperateTime+5*60 < time.Now().Unix() {
		response.ThrowException(c, exception.ExceptionInvalidOperation)
		return
	}

	// check the sign
	verifySign := map[string]string{}
	verifySign["operator"] = req.Operator
	verifySign["operateTime"] = string(req.OperateTime)
	verifySign["clientIP"] = req.ClientIP

	//if !util.VerifySIGN(verifySign, req.Sign) {
	//	response.ThrowException(c, exception.ExceptionIllegalSign)
	//	return
	//}

	zap.L().Info("GetRoamMiniWxAPIDataHandler", zap.String("operator", req.Operator), zap.String("ip", req.ClientIP))

	miniToken, err := interior.GetAccessToken()

	if err != nil {
		zap.L().Error("GetRoamMiniWxAPIDataHandler.GetMiniAccessTokenError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetMiniAccTokenFailure)
		return
	}

	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("20060102")

	yVisitTotal, err := interior.GetAnalyDailySummary(miniToken, yesterdayStr, yesterdayStr)

	if err != nil {
		zap.L().Error("GetRoamMiniWxAPIDataHandler.GetAnalyDailySummaryError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetDailSummaryFailure)
		return
	}

	visitPv, visitUv, visitUvNew, stayTimeUv, err := interior.GetAnalyDailyVisitTrend(miniToken, yesterdayStr, yesterdayStr)

	if err != nil {
		zap.L().Error("GetRoamMiniWxAPIDataHandler.GetAnalyDailyVisitTrend", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetVisitTrendFailure)
		return
	}

	rpageVisitPv, rpageVisitUv, dpageVisitPv, dpageVisitUv, err := interior.GetAnalyVisitPage(miniToken, yesterdayStr, yesterdayStr)

	if err != nil {
		zap.L().Error("GetRoamMiniWxAPIDataHandler.GetAnalyVisitPage", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetVisitPageFailure)
		return
	}

	resp := &protocol.GetMiniDataResp{
		MiniVisitTotal:     yVisitTotal,
		MiniVisitPv:        visitPv,
		MiniVisitUv:        visitUv,
		MiniVisitUvNew:     visitUvNew,
		MiniStayTimeUv:     stayTimeUv,
		MiniRcyPageVisitPv: rpageVisitPv,
		MiniRcyPageVisitUv: rpageVisitUv,
		MiniDetPageVisitPv: dpageVisitPv,
		MiniDetPageVisitUv: dpageVisitUv,
	}

	response.SuccessResp(c, resp)
}
