package timed

import (
	"github.com/robfig/cron"
)

func SetupTimed() error {
	crn := cron.New()

	var err error

	// 凌晨3点执行获取用户总览数据任务
	if err = crn.AddFunc("0 0 3 * * ? ", GetOverViewData); err != nil {
		return err
	}

	// 早上9点执行获取微信公众号数据任务
	if err = crn.AddFunc("0 0 9 * * ? ", GeAppletsWxAPIData); err != nil {
		return err
	}

	// 早上9点10分执行获取微信小程序数据任务
	if err = crn.AddFunc("0 10 9 * * ? ", GetMiniWxAPIData); err != nil {
		return err
	}

	// 早上9点20分执行获取销售订单数据任务
	if err = crn.AddFunc("0 20 9 * * ? ", GetOrderStaticticalData); err != nil {
		return err
	}

	// 早上9点30分执行获取回收订单数据任务
	if err = crn.AddFunc("0 30 9 * * ? ", GetRcyOrderStatisticalData); err != nil {
		return err
	}

	// 凌晨1点执行一次补充到货通知openid
	if err = crn.AddFunc("0 0 1 * * ? ", UpdateArrivalNoticeOpenid); err != nil {
		return err
	}

	// 凌晨1点30分再次执行一次补充到货通知openid
	if err = crn.AddFunc("0 30 1 * * ? ", UpdateArrivalNoticeOpenid); err != nil {
		return err
	}

	// 凌晨2点再一次执行补充到货通知openid
	if err = crn.AddFunc("0 0 2 * * ? ", UpdateArrivalNoticeOpenid); err != nil {
		return err
	}

	// 凌晨3点30分执行一次补充到货通知unionid
	if err = crn.AddFunc("0 30 3 * * ? ", SupplementaryUnionid); err != nil {
		return err
	}

	// 凌晨4点再执行执行一次补充到货通知unionid
	if err = crn.AddFunc("0 0 4 * * ? ", SupplementaryUnionid); err != nil {
		return err
	}

	// 凌晨4点30分执行一次补充到货通知unionid
	if err = crn.AddFunc("0 30 4 * * ? ", SupplementaryUnionid); err != nil {
		return err
	}

	// 早上5点执行一次补充到货通知openid_wa
	if err = crn.AddFunc("0 0 5 * * ? ", UpdateArrivalNoticeOpenidWa); err != nil {
		return err
	}

	// 早上5点20分再次执行补充到货通知openid_wa
	if err = crn.AddFunc("0 20 5 * * ? ", UpdateArrivalNoticeOpenidWa); err != nil {
		return err
	}

	// 早上5点40分再次执行补充到货通知openid_wa
	if err = crn.AddFunc("0 40 5 * * ? ", UpdateArrivalNoticeOpenidWa); err != nil {
		return err
	}

	crn.Start()

	return nil
}
