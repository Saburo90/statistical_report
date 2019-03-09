package timed

import (
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/models/arrival_model"
	"gitee.com/NotOnlyBooks/statistical_report/models/users_model"
	"go.uber.org/zap"
	"time"
)

func UpdateArrivalNoticeOpenid() {
	session := components.NewDBSession()

	defer session.Close()

	arrivalNoticeModel := arrival_model.NewArrivalNoticeModel(session)

	memberModel := users_model.NewEweiShopMemberModel(session)

	zap.L().Info("开始获取到货通知表中openid为空数据(限制20000条)")

	var (
		emptyOpenidArri []arrival_model.ArrivalNotice
		mOpenid         string
		err             error
	)
	for retry := 3; retry > 0; retry-- {
		emptyOpenidArri, err = arrivalNoticeModel.GetEmptyOpenidArrival(uint64(constant.RoamingApplets))

		if err != nil {
			zap.L().Error("获取openid为空的到货通知失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取openid为空的到货通知失败, 重试机会已耗尽, 此处更新空openid到货通知失败")
			return
		}
	}
	if len(emptyOpenidArri) > 0 {
		successNum := 0
		failureNum := 0
		for _, arrival := range emptyOpenidArri {
			// 获取用户openid
			mOpenid, err = memberModel.GetMemberOpenid(arrival.ArriUnionid)

			if err == nil && mOpenid != "" {
				uarriBean := map[string]interface{}{}
				uarriBean["openid"] = mOpenid
				//更新openid
				update_res := arrivalNoticeModel.UpdateOpenid(uarriBean, arrival.Id)

				if !update_res {
					failureNum++
					zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "更新openid时故障"))
				} else {
					successNum++
				}
			} else {
				failureNum++
				zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "获取会员表中openid失败"))
			}
		}

		zap.L().Info("本轮更新结束", zap.Int("成功", successNum), zap.Int("失败", failureNum))
	}

}
