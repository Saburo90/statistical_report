package sale_statis

import (
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/exception"
	"github.com/Saburo90/statistical_report/models/sales_model"
	"github.com/Saburo90/statistical_report/models/statis_model"
	"github.com/Saburo90/statistical_report/models/users_model"
	"github.com/Saburo90/statistical_report/protocol"
	"github.com/Saburo90/statistical_report/request"
	"github.com/Saburo90/statistical_report/response"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"time"
)

type SaleReq struct {
	Operator    string `json:"operator"`
	OperateTime int64  `json:"operateTime"`
	ClientIP    string `json:"clientIP"`
	Sign        string `json:"sign"`
}

func (req *SaleReq) Validation() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Operator, validation.Required, validation.In("admin", "saburo", "system")),
		validation.Field(&req.OperateTime, validation.Required),
		validation.Field(&req.ClientIP, validation.Required, validation.Match(regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"))),
		validation.Field(&req.Sign, validation.Required),
	)
}

func GetSalesDataHandler(c *gin.Context) {
	req := &SaleReq{}

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

	session := components.NewDBSession()

	defer session.Close()

	eweiShopOrderModel := sales_model.NewEweiShopOrderModel(session)

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1

	yNOrder, yNOrderTPrice, err := eweiShopOrderModel.GetOrderStatisData(constant.RoamingApplets, yesterdayStart, yesterdayEnd, -1)

	if err != nil {
		zap.L().Error("GetSalesDataHandler.GetOrderDataError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetOrderDataFailure)
		return
	}

	meanOrdPrice := decimal.NewFromFloat(yNOrderTPrice).Div(decimal.New(yNOrder, 0)).Round(4)
	conversionOrdRate := decimal.New(0, 0)
	visitUv, err := components.Redis.Get(constant.RedisYesNewVisitUv).Result()

	if err != nil || err == redis.Nil {
		statis := &statis_model.StatisticalReport{}

		exist, err := session.Select("total_rmini_newvistsuv").Get(statis)

		if err != nil {
			zap.L().Error("获取昨日小程序访问人数失败", zap.Error(err))
			response.ThrowException(c, exception.ExceptionGetVisitUvFailure)
			return
		}

		if !exist {
			response.ThrowException(c, exception.ExceptionGetVisitUvFailure)
			return
		}
		if statis.TotalRminiNewVisitsUv != 0 {
			conversionOrdRate = decimal.New(yNOrder, 0).Div(decimal.New(statis.TotalRminiNewVisitsUv, 0)).Round(4)
		} else {
			response.ThrowException(c, exception.ExceptionGetVisitUvFailure)
			return
		}
	} else {
		new_VisitUv, err := strconv.ParseInt(visitUv, 10, 64)
		if err != nil {
			zap.L().Error("转化昨日小程序访问人数失败", zap.Error(err))
			response.ThrowException(c, exception.ExceptionSystemError)
			return
		}
		conversionOrdRate = decimal.New(yNOrder, 0).Div(decimal.New(new_VisitUv, 0)).Round(4)
	}

	yNOrderGtotal, err := session.Table("ims_ewei_shop_order_goods").
		Join("LEFT", "ims_ewei_shop_order", "ims_ewei_shop_order_goods.orderid = ims_ewei_shop_order.id").
		Where("ims_ewei_shop_order_goods.createtime > ?", yesterdayStart).
		And("ims_ewei_shop_order_goods.createtime <= ?", yesterdayEnd).
		And("ims_ewei_shop_order.createtime > ?", yesterdayStart).
		And("ims_ewei_shop_order.createtime <= ?", yesterdayEnd).
		And("ims_ewei_shop_order.status != ?", -1).Count(&sales_model.OrderGoods{})

	if err != nil {
		zap.L().Error("GetSalesDataHandler.GetYNewOrderGoodsError", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetOrderGoodsTotalFailure)
		return
	}

	meamOrdGoods := decimal.New(yNOrderGtotal, 0).Div(decimal.New(yNOrder, 0)).Round(4)

	shopMemModel := users_model.NewEweiShopMemberModel(session)

	afterTomorrorEnd := yesterdayStart - 1
	aftStr := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	aft, _ := time.ParseInLocation("2006-01-02", aftStr, time.Local)
	afterTomorrorStart := aft.Unix()

	alterTomorrorU, err := shopMemModel.GetTimeRangeUsers(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd)

	if err != nil {
		zap.L().Error("获取注册时间为前天范围内用户数据失败", zap.Error(err))
		response.ThrowException(c, exception.ExceptionGetRangeTimeUsersFailure)
		return
	}

	var (
		unionidStr       []string
		orderNum24       int64
		ordRateIn24      decimal.Decimal
		aftStatis        statis_model.StatisticalReport
		orderPirce24     decimal.Decimal
		meanOrdPrice24   decimal.Decimal
		ordGoodsIn24     int64
		meamOrdGoodsIn24 decimal.Decimal
		ordNumErr24      error
	)

	if len(alterTomorrorU) > 0 {
		for _, rangeU := range alterTomorrorU {
			unionidStr = append(unionidStr, rangeU.Unionid)
		}
		if len(unionidStr) > 0 {
			orderNum24, ordNumErr24 = eweiShopOrderModel.CountRangeOrder(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd, unionidStr)

			if ordNumErr24 != nil {
				zap.L().Error("获取24小时销售下单数失败", zap.Error(err))
				response.ThrowException(c, exception.ExceptionGetOrderNumIn24Failure)
				return
			}

			if orderNum24 > 0 {
				exist, err := session.Select("total_rmini_newvistusr").Where("daily = ?", aftStr).Get(&aftStatis)

				if err != nil || !exist {
					ordRateIn24 = decimal.New(orderNum24, 0).Div(decimal.New(int64(len(unionidStr)), 0)).Round(4)
				} else {
					ordRateIn24 = decimal.New(orderNum24, 0).Div(decimal.New(aftStatis.TotalRminiNewVisitsUsr, 0)).Round(4)
				}

				orderPirce24, ordNumErr24 = eweiShopOrderModel.CountRangeOrdersPrice(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd, unionidStr)

				if ordNumErr24 != nil {
					zap.L().Error("获取24小时销售订单总价失败", zap.Error(err))
					response.ThrowException(c, exception.ExceptionGetOrderPirceIn24Failure)
					return
				}

				if orderPirce24.Cmp(decimal.New(0, 0)) == 1 {
					meanOrdPrice24 = orderPirce24.Div(decimal.New(orderNum24, 0)).Round(4)
				}

				ordGoodsIn24, ordNumErr24 = session.Table("ims_ewei_shop_order_goods").
					Join("LEFT", "ims_ewei_shop_order", "ims_ewei_shop_order_goods.orderid = ims_ewei_shop_order.id").
					Where("ims_ewei_shop_order_goods.uniacid = ?", constant.RoamingApplets).
					And("ims_ewei_shop_order.uniacid = ?", constant.RoamingApplets).
					And("ims_ewei_shop_order.status != ?", -1).
					In("ims_ewei_shop_order.unionid", unionidStr).
					Count(&sales_model.OrderGoods{})

				if ordNumErr24 != nil {
					zap.L().Error("获取24小时销售订单总商品数失败", zap.Error(err))
					response.ThrowException(c, exception.ExceptionGetOrderPirceIn24Failure)
					return
				}

				meamOrdGoodsIn24 = decimal.New(ordGoodsIn24, 0).Div(decimal.New(orderNum24, 0)).Round(4)
			}
		}
	}

	resp := &protocol.GetOrderDataResp{
		YNorder:           yNOrder,
		YNorderPrice:      yNOrderTPrice,
		MeanOrdPrice:      meanOrdPrice,
		MeanOrdGoos:       meamOrdGoods,
		ConversionOrdRate: conversionOrdRate,
		TNorderGoodsT:     yNOrderGtotal,
		OrderNumIn24:      orderNum24,
		ConverOrdRate24:   ordRateIn24,
		MeanOrdPrice24:    meanOrdPrice24,
		MeanOrdGoodsIn24:  meamOrdGoodsIn24,
	}

	response.SuccessResp(c, resp)
}
