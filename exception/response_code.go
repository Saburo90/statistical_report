package exception

// System Error
const (
	SuccessCode = iota
	IllegalSignCode
	SystemErrorCode
	BusyServerCode
	IllegalParameterCode
	ResultRunAwayCode
	IllegalOperateCode
	SystemMaintenanceCode
	InvalidOperationCode
)

var (
	Success                    = New(SuccessCode, "OK")
	ExceptionIllegalSign       = New(IllegalSignCode, "ILLEGAL SIGN")
	ExceptionSystemError       = New(SystemErrorCode, "SYSTEM ERROR")
	ExceptionBusyServer        = New(BusyServerCode, "BUSY SERVER")
	ExceptionIllegalParameter  = New(IllegalParameterCode, "ILLEGAL PARAMETER")
	ExceptionResultRunAway     = New(ResultRunAwayCode, "Oh, No! The Result Run Away")
	ExceptionIllegalOperate    = New(IllegalOperateCode, "ILLEGAL OPERATE")
	ExceptionSystemMaintenance = New(SystemMaintenanceCode, "SYSTEM MAINTENANCE")
	ExceptionInvalidOperation  = New(InvalidOperationCode, "INVALID OPERATION")
)

// GetWxAPIData
const (
	GetRoamAppletsAccessTokenFailure = 1001
	GetCumulateUsersFailure          = 1002
	GetSummaryUsersFailure           = 1003
	GetOrderDataFailure              = 1004
	GetMiniAccTokenFailure           = 1005
	GetOrderGoodsTotalFailure        = 1006
	GetDailSummaryFailure            = 1007
	GetVisitTrendFailure             = 1008
	GetVisitPageFailure              = 1009
	GetVisitUvFailure                = 1010
	GetRangeTimeUsersFailure         = 1011
	GetOrderNumIn24Failure           = 1012
	GetAftVisitUvNewUserFailure      = 1013
	NotFoundAftVisitUvNewUsers       = 1014
	GetOrderPirceIn24Failure         = 1015
	GetOrderGoodsIn24Failure         = 1015
)

var (
	ExceptionGetRoamAppletsAccessToken   = New(GetRoamAppletsAccessTokenFailure, "GET ROAM APPLETS ACCESS TOKEN FAILURE")
	ExceptionGetCumulateUsers            = New(GetCumulateUsersFailure, "GET CUMULATE USERS FAILURE")
	ExceptionGetSummaryUsers             = New(GetSummaryUsersFailure, "GET SUMMARY USERS FAILURE")
	ExceptionGetOrderDataFailure         = New(GetOrderDataFailure, "GET ORDER DATA FAILURE")
	ExceptionGetMiniAccTokenFailure      = New(GetMiniAccTokenFailure, "GET MINI ACCESS TOKEN FAILURE")
	ExceptionGetOrderGoodsTotalFailure   = New(GetOrderGoodsTotalFailure, "GET ORDER GOODS TOTAL FAILURE")
	ExceptionGetDailSummaryFailure       = New(GetDailSummaryFailure, "GET DAILY SUMMARY FAILURE")
	ExceptionGetVisitTrendFailure        = New(GetVisitTrendFailure, "GET VISIT TREND FAILURE")
	ExceptionGetVisitPageFailure         = New(GetVisitPageFailure, "GET VISIT PAGE FAILURE")
	ExceptionGetVisitUvFailure           = New(GetVisitUvFailure, "GET VISIT UV FAILURE")
	ExceptionGetRangeTimeUsersFailure    = New(GetRangeTimeUsersFailure, "GET RANGETIME USERS FAILURE")
	ExceptionGetOrderNumIn24Failure      = New(GetOrderNumIn24Failure, "GET ORDER NUM IN 24 HOURS FAILURE")
	ExceptionGetAftVisitUvNewUserFailure = New(GetAftVisitUvNewUserFailure, "GET AFT VISIT UV NEW USERS")
	ExceptionNotFoundAftVisitUvNewUsers  = New(NotFoundAftVisitUvNewUsers, "NOT FOUND VISIT UV NEW USERS")
	ExceptionGetOrderPirceIn24Failure    = New(GetOrderPirceIn24Failure, "GET ORDER PRICE IN 24 HOURS FAILURE")
	ExceptionGetOrderGoodsIn24Failure    = New(GetOrderGoodsIn24Failure, "GET ORDER PRICE IN 24 HOURS FAILURE")
)
