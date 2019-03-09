package timed

import (
	"fmt"
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/models/sales_model"
	"github.com/Saburo90/statistical_report/models/statis_model"
	"github.com/Saburo90/statistical_report/models/users_model"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GetOrderStaticticalData() {

	session := components.NewDBSession()

	defer session.Close()

	eweiShopOrderModel := sales_model.NewEweiShopOrderModel(session)

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	var (
		yNOrder           int64
		yNOrderTPrice     float64
		meanOrdPrice      decimal.NullDecimal
		conversionOrdRate decimal.NullDecimal
		yNOrderGtotal     int64
		meamOrdGoods      decimal.NullDecimal
		alterTomorrorU    []users_model.Member
		err               error
		unionidStr        []string
		orderNum24        int64
		ordRateIn24       decimal.NullDecimal
		aftStatis         statis_model.GetAftStatisRecord
		orderPirce24      decimal.Decimal
		meanOrdPrice24    decimal.NullDecimal
		ordGoodsIn24      int64
		meamOrdGoodsIn24  decimal.NullDecimal
		dailyCount        int64
		ordNumErr24       error
	)

	meanOrdPrice.Decimal = decimal.New(0, 0)
	conversionOrdRate.Decimal = decimal.New(0, 0)
	meamOrdGoods.Decimal = decimal.New(0, 0)
	ordRateIn24.Decimal = decimal.New(0, 0)
	meanOrdPrice24.Decimal = decimal.New(0, 0)
	meamOrdGoodsIn24.Decimal = decimal.New(0, 0)
	meanOrdPrice.Valid = true
	conversionOrdRate.Valid = true
	meamOrdGoods.Valid = true
	ordRateIn24.Valid = true
	meanOrdPrice24.Valid = true
	meamOrdGoodsIn24.Valid = true

	zap.L().Info("开始获取昨日销售订单数")
	for retry := 3; retry > 0; retry-- {
		yNOrder, yNOrderTPrice, err = eweiShopOrderModel.GetOrderStatisData(constant.RoamingApplets, yesterdayStart, yesterdayEnd, -1)

		if err != nil {
			zap.L().Error("获取漫游鲸昨日销售(订单数/订单总价)失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸昨日销售(订单数/订单总价)失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
			return
		}
	}

	zap.L().Info("开始获取销售转化率")
	for retry := 5; retry > 0; retry-- {
		visitUv, err := components.Redis.Get(constant.RedisYesNewVisitUv).Result()

		if err != nil || err == redis.Nil {
			var statis statis_model.GetYesterdayRecord

			exist, err := session.Select("total_rmini_newvistsuv").
				Table("ims_statistical_report").
				Where("daily = ?", yesterdayStr).
				Get(&statis)

			if err != nil {
				zap.L().Error("获取漫游鲸前天小程序新增用户数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(time.Second * 3)
			}

			if !exist {
				zap.L().Error("获取漫游鲸前天小程序新增用户数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(time.Second * 3)
			}
			if statis.TotalRminiNewVisitsUv != 0 {
				conversionOrdRate.Decimal = decimal.New(yNOrder, 0).Div(decimal.New(statis.TotalRminiNewVisitsUv, 0)).Round(4)
				conversionOrdRate.Valid = true
				break
			} else {
				zap.L().Error("漫游鲸前天小程序新增用户数为0, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(time.Second * 3)
			}
		} else {
			new_VisitUv, err := strconv.ParseInt(visitUv, 10, 64)
			if err != nil {
				zap.L().Error("转化昨日小程序访问人数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(time.Second * 3)
			}
			conversionOrdRate.Decimal = decimal.New(yNOrder, 0).Div(decimal.New(new_VisitUv, 0)).Round(4)
			conversionOrdRate.Valid = true
			break
		}

		if retry == 0 {
			zap.L().Error("获取销售转化率失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
			return
		}
	}

	zap.L().Info("开始获取昨日销售订单商品总数")
	for retry := 3; retry > 0; retry-- {
		yNOrderGtotal, err = session.Table("ims_ewei_shop_order_goods").
			Join("LEFT", "ims_ewei_shop_order", "ims_ewei_shop_order_goods.orderid = ims_ewei_shop_order.id").
			Where("ims_ewei_shop_order_goods.createtime > ?", yesterdayStart).
			And("ims_ewei_shop_order_goods.createtime <= ?", yesterdayEnd).
			And("ims_ewei_shop_order.createtime > ?", yesterdayStart).
			And("ims_ewei_shop_order.createtime <= ?", yesterdayEnd).
			And("ims_ewei_shop_order.status != ?", -1).Count(&sales_model.OrderGoods{})

		if err != nil {
			zap.L().Error("获取昨日销售订单商品总数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取昨日销售订单商品总数失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
			return
		}
	}

	if yNOrder > 0 {
		zap.L().Info("开始计算销售客单价")
		meanOrdPrice.Decimal = decimal.NewFromFloat(yNOrderTPrice).Div(decimal.New(yNOrder, 0)).Round(2)
		meanOrdPrice.Valid = true
		zap.L().Info("开始计算单均销量")
		meamOrdGoods.Decimal = decimal.New(yNOrderGtotal, 0).Div(decimal.New(yNOrder, 0)).Round(2)
		meamOrdGoods.Valid = true
	}

	shopMemModel := users_model.NewEweiShopMemberModel(session)

	afterTomorrorEnd := yesterdayStart - 1
	aftStr := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	aft, _ := time.ParseInLocation("2006-01-02", aftStr, time.Local)
	afterTomorrorStart := aft.Unix()

	zap.L().Info("开始获取前天注册用户数")
	for retry := 3; retry > 0; retry-- {
		alterTomorrorU, err = shopMemModel.GetTimeRangeUsers(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd)

		if err != nil {
			zap.L().Error("获取前天注册用户数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取前天注册用户数失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
			return
		}
	}

	if len(alterTomorrorU) > 0 {
		for _, rangeU := range alterTomorrorU {
			unionidStr = append(unionidStr, rangeU.Unionid)
		}
		if len(unionidStr) > 0 {
			zap.L().Info("开始获取24小时销售订单数")
			for retry := 3; retry > 0; retry-- {
				orderNum24, ordNumErr24 = eweiShopOrderModel.CountRangeOrder(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd, unionidStr)

				if ordNumErr24 != nil {
					zap.L().Error("获取24小时销售订单数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
					time.Sleep(time.Second * 3)
				} else {
					break
				}

				if retry == 0 {
					zap.L().Error("获取24小时销售订单数失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
					return
				}
			}

			if orderNum24 > 0 {
				zap.L().Info("开始计算24小时销售转化率")
				exist, err := session.Select("total_rmini_newvistusr").
					Table("ims_statistical_report").
					Where("daily = ?", aftStr).
					Get(&aftStatis)

				if err != nil || !exist {
					ordRateIn24.Decimal = decimal.New(orderNum24, 0).Div(decimal.New(int64(len(unionidStr)), 0)).Round(4)
					ordRateIn24.Valid = true
				} else {
					ordRateIn24.Decimal = decimal.New(orderNum24, 0).Div(decimal.New(aftStatis.TotalRminiNewVisitsUsr, 0)).Round(4)
					ordRateIn24.Valid = true
				}

				zap.L().Info("开始获取24小时销售订单总价")
				for retry := 3; retry > 0; retry-- {
					orderPirce24, ordNumErr24 = eweiShopOrderModel.CountRangeOrdersPrice(constant.RoamingApplets, afterTomorrorStart, afterTomorrorEnd, unionidStr)

					if ordNumErr24 != nil {
						zap.L().Error("获取24小时销售订单总价失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
						time.Sleep(time.Second * 3)
					} else {
						break
					}

					if retry == 0 {
						zap.L().Error("获取24小时销售订单总价失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
						return
					}
				}

				if orderPirce24.Cmp(decimal.New(0, 0)) == 1 {
					zap.L().Info("开始计算24销售客单价")
					meanOrdPrice24.Decimal = orderPirce24.Div(decimal.New(orderNum24, 0)).Round(2)
					meanOrdPrice24.Valid = true
				}

				zap.L().Info("开始获取24小时销售订单商品总数")
				for retry := 3; retry > 0; retry-- {
					ordGoodsIn24, ordNumErr24 = session.Table("ims_ewei_shop_order_goods").
						Join("LEFT", "ims_ewei_shop_order", "ims_ewei_shop_order_goods.orderid = ims_ewei_shop_order.id").
						Where("ims_ewei_shop_order_goods.uniacid = ?", constant.RoamingApplets).
						And("ims_ewei_shop_order.uniacid = ?", constant.RoamingApplets).
						And("ims_ewei_shop_order_goods.createtime > ?", afterTomorrorStart).
						And("ims_ewei_shop_order_goods.createtime <= ?", yesterdayEnd).
						And("ims_ewei_shop_order.createtime > ?", afterTomorrorStart).
						And("ims_ewei_shop_order.createtime <= ?", yesterdayEnd).
						And("ims_ewei_shop_order.status != ?", -1).
						In("ims_ewei_shop_order.unionid", unionidStr).
						Count(&sales_model.OrderGoods{})

					if ordNumErr24 != nil {
						zap.L().Error("获取24小时销售订单商品总数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
						time.Sleep(time.Second * 3)
					} else {
						break
					}

					if retry == 0 {
						zap.L().Error("获取24小时销售订单商品总数失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
						return
					}
				}
				zap.L().Info("开始计算24小时单均销量")
				meamOrdGoodsIn24.Decimal = decimal.New(ordGoodsIn24, 0).Div(decimal.New(orderNum24, 0)).Round(2)
				meamOrdGoodsIn24.Valid = true
			}

		}
	}

	zap.L().Info("开始检验是否存在昨日统计记录")
	staticticalModel := statis_model.NewStatisticalReportModel(session)

	for retry := 3; retry > 0; retry-- {

		dailyCount, err = staticticalModel.CheckDaily(yesterdayStr)

		if err != nil {
			zap.L().Error("检验是否存在昨日统计记录失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("检验是否存在昨日统计记录失败，重试机会已耗尽，强制新增记录")
			dailyCount = 0
			break
		}
	}

	if dailyCount == 0 {
		// 新增记录
		zap.L().Info("开始记录销售订单数据")
		orderStatis := statis_model.OrderStatis{
			Daily:                 yesterdayStr,
			TotalRorderNewNo:      yNOrder,
			TotalRorderGoodsNewNo: yNOrderGtotal,
			MeanOrdPrice:          meanOrdPrice,
			MeanOrdGoods:          meamOrdGoods,
			ConversionOrdRate:     conversionOrdRate,
			OrderNumIn24:          orderNum24,
			ConversionOrdRate24:   ordRateIn24,
			MeanOrdPrice24:        meanOrdPrice24,
			MeanOrdGoods24:        meamOrdGoodsIn24,
		}
		for retry := 5; retry > 0; retry-- {
			_, err = session.Table("ims_statistical_report").
				InsertOne(&orderStatis)

			if err != nil {
				fmt.Println("err = ", err)
				zap.L().Error("记录销售订单数据数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("销售订单数据记录成功", zap.String("记录日历", yesterdayStr), zap.String("记录时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("记录销售订单数据失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
				return
			}
		}

	} else {
		// 更新记录
		zap.L().Info("开始更新销售订单数据")

		for retry := 5; retry > 0; retry-- {
			orderStatisUpd := statis_model.OrderStatis{
				Daily:                 yesterdayStr,
				TotalRorderNewNo:      yNOrder,
				TotalRorderGoodsNewNo: yNOrderGtotal,
				MeanOrdPrice:          meanOrdPrice,
				MeanOrdGoods:          meamOrdGoods,
				ConversionOrdRate:     conversionOrdRate,
				OrderNumIn24:          orderNum24,
				ConversionOrdRate24:   ordRateIn24,
				MeanOrdPrice24:        meanOrdPrice24,
				MeanOrdGoods24:        meamOrdGoodsIn24,
			}

			_, err := session.Table("ims_statistical_report").
				Where("daily = ?", yesterdayStr).
				Update(&orderStatisUpd)

			if err != nil {
				zap.L().Error("更新销售订单数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("销售订单数据更新成功", zap.String("更新日历", yesterdayStr), zap.String("更新时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("更新销售订单数据失败，重试机会已耗尽，今日获取漫游鲸销售订单数据失败")
				return
			}
		}
	}

}
