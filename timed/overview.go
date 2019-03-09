package timed

import (
	"fmt"
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/models/statis_model"
	"github.com/Saburo90/statistical_report/models/users_model"
	"go.uber.org/zap"
	"time"
)

func GetOverViewData() {
	session := components.NewDBSession()

	defer session.Close()

	eweiShopMemMode := users_model.NewEweiShopMemberModel(session)

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1

	var (
		tPUsers    int64
		tRusers    int64
		tRMusers   int64
		tPNusers   int64
		tRNusers   int64
		tSNusers   int64
		tRANusers  int64
		tRMNusers  int64
		dailyCount int64
		err        error
	)

	zap.L().Info("开始获取漫游鲸总用户数")
	for retry := 3; retry > 0; retry-- {
		tPUsers, err = eweiShopMemMode.GetPlatRUsersNum(constant.RoamingApplets, yesterdayEnd)

		if err != nil {
			zap.L().Error("获取漫游鲸总用户数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸总用户数失败，重试机会已耗尽, 今日用户总览数据获取失败")
			return
		}
	}

	zap.L().Info("开始获取纯漫游鲸公众号用户数")
	for retry := 3; retry > 0; retry-- {
		tRusers, err = eweiShopMemMode.GetRoamAppletsUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

		if err != nil {
			zap.L().Error("获取纯漫游鲸公众号用户数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取纯漫游鲸公众号用户数失败，重试机会已耗尽，今日用户总览数据获取失败")
			return
		}
	}

	zap.L().Info("开始获取纯漫游鲸小程序用户数")
	for retry := 3; retry > 0; retry-- {
		tRMusers, err = eweiShopMemMode.GetRoamMiniUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

		if err != nil {
			zap.L().Error("获取纯漫游鲸小程序用户数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取纯漫游鲸小程序用户数失败，重试机会已耗尽，今日用户总览数据获取失败")
			return
		}
	}

	zap.L().Info("开始获取(漫游鲸/平台/共享店)昨日新增用户数")
	for retry := 3; retry > 0; retry-- {
		tRNusers, tPNusers, tSNusers, err = eweiShopMemMode.GetYesterdayPNewUsersNum(constant.RoamingApplets, constant.ShareApplets, "", yesterdayStart, yesterdayEnd)

		if err != nil {
			zap.L().Error("获取(漫游鲸/平台/共享店)昨日新增用户数，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取(漫游鲸/平台/共享店)昨日新增用户数，重试机会已耗尽，今日用户总览数据获取失败")
			return
		}
	}

	zap.L().Info("开始获取漫游鲸纯(公众号/小程序)昨日新增用户数")
	for retry := 3; retry > 0; retry-- {
		tRANusers, tRMNusers, err = eweiShopMemMode.GetRapltAndRminiNewUsersNum(constant.RoamingApplets, "", "", yesterdayStart, yesterdayEnd)

		if err != nil {
			zap.L().Error("获取漫游鲸纯(公众号/小程序)昨日新增用户数数，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸纯(公众号/小程序)昨日新增用户数，重试机会已耗尽，今日用户总览数据获取失败")
			return
		}
	}

	staticticalModel := statis_model.NewStatisticalReportModel(session)

	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	zap.L().Info("开始检验是否存在昨日统计记录")
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
		zap.L().Info("开始记录用户总览统计数据")
		overView := statis_model.OverViewStatis{
			Daily:           yesterdayStr,
			TotalRoamU:      tPUsers,
			TotalRapltsU:    tRusers,
			TotalRminiU:     tRMusers,
			TotalRoamNewU:   tRNusers,
			TotalRplatNewU:  tPNusers,
			TotalSapltsNewU: tSNusers,
			TotalRapltsNewU: tRANusers,
			TotalRminiNewU:  tRMNusers,
		}
		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				InsertOne(&overView)

			if err != nil {
				fmt.Println("err = ", err)
				zap.L().Error("记录用户总览数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("用户总览数据记录成功", zap.String("记录日历", yesterdayStr), zap.String("记录时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("记录用户总览数据失败，重试机会已耗尽，今日用户总览数据获取失败")
				return
			}
		}

	} else {
		// 更新记录
		zap.L().Info("开始更新用户总览统计数据")
		overViewUpd := statis_model.OverViewStatis{
			Daily:           yesterdayStr,
			TotalRoamU:      tPUsers,
			TotalRapltsU:    tRusers,
			TotalRminiU:     tRMusers,
			TotalRoamNewU:   tRNusers,
			TotalRplatNewU:  tPNusers,
			TotalSapltsNewU: tSNusers,
			TotalRapltsNewU: tRANusers,
			TotalRminiNewU:  tRMNusers,
		}

		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				Where("daily = ?", yesterdayStr).
				Update(&overViewUpd)

			if err != nil {
				zap.L().Error("更新用户总览数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("用户总览数据更新成功", zap.String("更新日历", yesterdayStr), zap.String("更新时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("更新用户总览数据失败，重试机会已耗尽，今日用户总览数据获取失败")
				return
			}
		}
	}
}
