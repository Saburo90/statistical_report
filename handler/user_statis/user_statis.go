package user_statis

import (
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/exception"
	"gitee.com/NotOnlyBooks/statistical_report/handler/user_statis/interior"
	"gitee.com/NotOnlyBooks/statistical_report/models/users_model"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"gitee.com/NotOnlyBooks/statistical_report/request"
	"gitee.com/NotOnlyBooks/statistical_report/response"
	//"gitee.com/NotOnlyBooks/statistical_report/util"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type Requ struct {
	Operator    string `json:"operator"`
	OperateTime int64  `json:"operateTime"`
	ClientIP    string `json:"clientIP"`
	Sign        string `json:"sign"`
}

func (req *Requ) Validation() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Operator, validation.Required, validation.In("admin", "saburo", "system")),
		validation.Field(&req.OperateTime, validation.Required),
		validation.Field(&req.ClientIP, validation.Required, validation.Match(regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"))),
		validation.Field(&req.Sign, validation.Required),
	)
}

// @Summary 获取用户统计数据
// @Tags 用户统计
// @Produce json
// @Param operationSign body protocol.OverviewReq true "调用者,调用时间,调用者IP,签名"
// @Success 200 {object} protocol.OverviewResp "{"code":0,"msg":"OK","data":{}}"
// @Failure 400 "{"code":exceptionCode,"msg":exceptionMsg,"data":{}}"
// @Router /user/getOverview [post]
func GetOverviewHandler(c *gin.Context) {

	req := &Requ{}

	if err := request.Bind(c, req); err != nil {
		//c.Error(err)
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
	zap.L().Info("GetOverviewHandler", zap.String("operator", req.Operator), zap.String("ip", req.ClientIP))
	session := components.NewDBSession()

	defer session.Close()

	eweiShopMemMode := users_model.NewEweiShopMemberModel(session)

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1

	tPUsers, err := eweiShopMemMode.GetPlatRUsersNum(constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		zap.L().Error("GetOverviewHandler.GetPaltUserError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionBusyServer)
		return
	}

	tRusers, err := eweiShopMemMode.GetRoamAppletsUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		zap.L().Error("GetOverviewHandler.GetRoamApltUserError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionBusyServer)
		return
	}

	tRMusers, err := eweiShopMemMode.GetRoamMiniUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		zap.L().Error("GetOverviewHandler.GetRoamMiniUserError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionBusyServer)
		return
	}

	tRNusers, tPNusers, tSNusers, err := eweiShopMemMode.GetYesterdayPNewUsersNum(constant.RoamingApplets, constant.ShareApplets, "", yesterdayStart, yesterdayEnd)

	if err != nil {
		zap.L().Error("GetOverviewHandler.GetPYNUserError.GetRYNUsersError.GetSYNUsersError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionBusyServer)
		return
	}

	tRANusers, tRMNusers, err := eweiShopMemMode.GetRapltAndRminiNewUsersNum(constant.RoamingApplets, "", "", yesterdayStart, yesterdayEnd)

	if err != nil {
		zap.L().Error("GetOverviewHandler.GetRYANUserError.GetRYMNUserError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionBusyServer)
		return
	}

	resp := &protocol.OverviewResp{
		TotalPUsers:   tPUsers,
		TotalRUsers:   tRusers,
		TotalRMUsers:  tRMusers,
		TotalPNUsers:  tPNusers,
		TotalRNUsers:  tRNusers,
		TotalSNUsers:  tSNusers,
		TotalRANUsers: tRANusers,
		TotalRMNUsers: tRMNusers,
	}

	response.SuccessResp(c, resp)
}

func GetWxAPIDataHandler(c *gin.Context) {
	req := &Requ{}

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
	zap.L().Info("GetOverviewHandler", zap.String("operator", req.Operator), zap.String("ip", req.ClientIP))
	// get roam applest access token
	roamaccToken, err := interior.GetAppletsAccessToken(constant.RoamingApplets)

	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetRoamAppletsAccessTokenError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetRoamAppletsAccessToken)
		return
	}

	// get fans total
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	yrCumulateFans, err := interior.GetUsersCumulate(roamaccToken, yesterdayStr, yesterdayStr)

	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetCumulateUsersError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetCumulateUsers)
		return
	}

	// get yesterday new fans and yesterday cancel fans

	yrNewFans, yrCancelFans, err := interior.GetUsersSummary(roamaccToken, yesterdayStr, yesterdayStr)
	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetSummaryUsersError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetSummaryUsers)
		return
	}

	// get roam applest access token
	shareaccToken, err := interior.GetAppletsAccessToken(constant.ShareApplets)

	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetRoamAppletsAccessTokenError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetRoamAppletsAccessToken)
		return
	}

	// get fans total
	ysCumulateFans, err := interior.GetUsersCumulate(shareaccToken, yesterdayStr, yesterdayStr)

	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetCumulateUsersError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetCumulateUsers)
		return
	}

	// get yesterday new fans and yesterday cancel fans

	ysNewFans, ysCancelFans, err := interior.GetUsersSummary(shareaccToken, yesterdayStr, yesterdayStr)
	if err != nil {
		zap.L().Error("GetWxAPIDataHandler.GetSummaryUsersError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetSummaryUsers)
		return
	}

	resp := &protocol.GetWxDataResp{
		YRtotalUsers:  yrCumulateFans,
		YRnewUsers:    yrNewFans,
		YRcancelUsers: yrCancelFans,
		YStotalUsers:  ysCumulateFans,
		YSnewUsers:    ysNewFans,
		YScancelUsers: ysCancelFans,
	}

	response.SuccessResp(c, resp)
}
