package timed

import (
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/models/arrival_model"
	"github.com/Saburo90/statistical_report/models/users_model"
	"github.com/Saburo90/statistical_report/modules/wxapi"
	"github.com/Saburo90/statistical_report/protocol"
	"go.uber.org/zap"
	"time"
)

/*
 * 清空请求Slice，传入的slice对象地址不变。
 * params:
 *   s: slice对象指针，类型为*[]protocol.OpenidList
 * return:
 *   无
 */
func SliceClear(s *[]protocol.OpenidList) {
	*s = (*s)[0:0]
}

// 自动补全到货通知表中缺失unionid记录
func SupplementaryUnionid() {
	session := components.NewDBSession()

	defer session.Close()

	arrivalNoticeModel := arrival_model.NewArrivalNoticeModel(session)

	memberModel := users_model.NewEweiShopMemberModel(session)

	zap.L().Info("开始获取到货通知表中unionid为空的数据(限制5000条)")

	var (
		emptyUnionidArrivals []arrival_model.ArrivalNotice
		mUnionid             string
		roamaccToken         string
		aUnionid             string
		err                  error
	)

	for retry := 3; retry > 0; retry-- {
		emptyUnionidArrivals, err = arrivalNoticeModel.GetEmptyUnionidArrival(uint64(constant.RoamingApplets))

		if err != nil {
			zap.L().Error("获取unionid为空的到货通知失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取unionid为空的到货通知失败, 重试机会已耗尽, 此处更新空openid到货通知失败")
			return
		}
	}

	total := len(emptyUnionidArrivals)
	if total > 0 {
		successNum := 0
		failureNum := 0
		for _, arrival := range emptyUnionidArrivals {
			// 获取用户unionid
			if arrival.ArriOpenidwa != "" && arrival.ArriOpenid == "" {
				mUnionid, err = memberModel.GetMemberUnionidByOpenidWa(arrival.ArriOpenidwa)
			} else if arrival.ArriOpenid != "" && arrival.ArriOpenidwa == "" {
				mUnionid, err = memberModel.GetMemberUnionidByOpenid(arrival.ArriOpenid)
			} else if arrival.ArriOpenid != "" && arrival.ArriOpenidwa != "" {
				mUnionid, err = memberModel.GetMemberUnionidByOpenidWa(arrival.ArriOpenidwa)
			} else {
				failureNum++
				// 删除该记录
				err = arrivalNoticeModel.DeleteRecordById(arrival.Id, arrival.ArriUniacid)
				zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "此条记录数据异常"), zap.Error(err))
				continue
			}

			if err == nil && mUnionid != "" {
				uarriBean := map[string]interface{}{}
				uarriBean["unionid"] = mUnionid
				//更新openid
				update_res := arrivalNoticeModel.UpdateUnionid(uarriBean, arrival.Id)

				if !update_res {
					failureNum++
					zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "更新unionid时故障"))
				} else {
					successNum++
				}
			} else {
				if arrival.ArriOpenid != "" {
					// 数据库中无法获取，调用微信API获取
					roamaccToken, err = getAppletsAccessToken(constant.RoamingApplets)

					if err != nil {
						failureNum++
						zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "获取access_token失败"))
						continue
					}

					aUnionid, err = wxapi.GetAppletUnionidByWxAPI(roamaccToken, arrival.ArriOpenid)

					if err == nil && aUnionid != "" {
						uarriBean := map[string]interface{}{}
						uarriBean["unionid"] = aUnionid
						//更新openid
						update_res := arrivalNoticeModel.UpdateUnionid(uarriBean, arrival.Id)

						if !update_res {
							failureNum++
							zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "更新unionid时故障"))
						} else {
							successNum++
						}
					} else {
						failureNum++
						// 删除该记录
						err = arrivalNoticeModel.DeleteRecordById(arrival.Id, arrival.ArriUniacid)
						zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "调用微信API获取unionid失败"), zap.Error(err))
					}

				} else {
					failureNum++
					// 删除该记录
					err = arrivalNoticeModel.DeleteRecordById(arrival.Id, arrival.ArriUniacid)
					zap.L().Error("此条记录更新失败", zap.Uint64("id为：", arrival.Id), zap.String("原因：", "此条记录openid为空"), zap.Error(err))
				}

			}

		}
		zap.L().Info("本轮更新结束", zap.Int("成功", successNum), zap.Int("失败", failureNum))
	} else {
		zap.L().Info("无须更新", zap.Int("成功", 0), zap.Int("失败", 0))
	}

}
