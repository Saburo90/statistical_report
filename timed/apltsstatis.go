package timed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/conf"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/models/statis_model"
	"github.com/Saburo90/statistical_report/protocol"
	"github.com/Saburo90/statistical_report/util"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GeAppletsWxAPIData() {
	var (
		roamaccToken   string
		shareaccToken  string
		yrCumulateFans int
		yrNewFans      int
		yrCancelFans   int
		ysCumulateFans int
		ysNewFans      int
		ysCancelFans   int
		dailyCount     int64
		err            error
	)

	zap.L().Info("开始获取漫游鲸公众号access_token")
	for retry := 5; retry > 0; retry-- {
		roamaccToken, err = getAppletsAccessToken(constant.RoamingApplets)

		if err != nil {
			zap.L().Error("获取漫游鲸公众号access_token失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("本次获取漫游鲸公众号access_token失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	zap.L().Info("开始获取漫游鲸公众号截止昨日终止时间总粉丝数")
	for retry := 3; retry > 0; retry-- {
		yrCumulateFans, err = getUsersCumulate(roamaccToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游鲸公众号截止昨日终止时间总粉丝数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸公众号截止昨日终止时间总粉丝数失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	zap.L().Info("开始获取漫游鲸公众号昨日(新增/取关)粉丝数")
	for retry := 3; retry > 0; retry-- {
		yrNewFans, yrCancelFans, err = getUsersSummary(roamaccToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游鲸公众号昨日(新增/取关)粉丝数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸公众号昨日(新增/取关)粉丝数失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	zap.L().Info("开始获取漫游共享店公众号access_token")
	for retry := 5; retry > 0; retry-- {
		shareaccToken, err = getAppletsAccessToken(constant.ShareApplets)

		if err != nil {
			zap.L().Error("获取漫游共享店公众号access_token失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("本次获取漫游共享店公众号access_token失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	zap.L().Info("开始获取漫游共享店公众号截止昨日终止时间总粉丝数")
	for retry := 3; retry > 0; retry-- {
		ysCumulateFans, err = getUsersCumulate(shareaccToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游共享店公众号截止昨日终止时间总粉丝数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游共享店公众号截止昨日终止时间总粉丝数失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	zap.L().Info("开始获取漫游共享店公众号昨日(新增/取关)粉丝数")
	for retry := 3; retry > 0; retry-- {
		ysNewFans, ysCancelFans, err = getUsersSummary(shareaccToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游共享店公众号昨日(新增/取关)粉丝数失败，3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游共享店公众号昨日(新增/取关)粉丝数失败，重试机会已耗尽，今日获取漫游鲸公众号粉丝数据失败")
			return
		}
	}

	zap.L().Info("开始检验是否存在昨日统计记录")
	session := components.NewDBSession()

	defer session.Close()

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
		zap.L().Info("开始记录漫游/共享店粉丝数据")
		apltsStatis := statis_model.AppletsStatis{
			Daily:            yesterdayStr,
			TotalRapltsFans:  int64(yrCumulateFans),
			TotalRapltsNewF:  int64(yrNewFans),
			TotalRapltsNewCF: int64(yrCancelFans),
			TotalSapltsFans:  int64(ysCumulateFans),
			TotalSapltsNewF:  int64(ysNewFans),
			TotalSapltsNewCF: int64(ysCancelFans),
		}
		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				InsertOne(&apltsStatis)

			if err != nil {
				fmt.Println("err = ", err)
				zap.L().Error("记录用户总览数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("漫游/共享店粉丝数据记录成功", zap.String("记录日历", yesterdayStr), zap.String("记录时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("记录漫游/共享店粉丝数据失败，重试机会已耗尽，今日漫游/共享店粉丝数据获取失败")
				return
			}
		}

	} else {
		// 更新记录
		zap.L().Info("开始更新漫游/共享店粉丝数据")
		apltsStatisUpd := statis_model.AppletsStatis{
			Daily:            yesterdayStr,
			TotalRapltsFans:  int64(yrCumulateFans),
			TotalRapltsNewF:  int64(yrNewFans),
			TotalRapltsNewCF: int64(yrCancelFans),
			TotalSapltsFans:  int64(ysCumulateFans),
			TotalSapltsNewF:  int64(ysNewFans),
			TotalSapltsNewCF: int64(ysCancelFans),
		}

		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				Where("daily = ?", yesterdayStr).
				Update(&apltsStatisUpd)

			if err != nil {
				zap.L().Error("更新漫游/共享店粉丝数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("漫游/共享店粉丝数据更新成功", zap.String("更新日历", yesterdayStr), zap.String("更新时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("更新漫游/共享店粉丝数据失败，重试机会已耗尽，今日用户总览数据获取失败")
				return
			}
		}
	}
}

func getAppletsAccessToken(uniacid uint64) (string, error) {
	// first to get access token from phpapi
	supplyTokenUrl := ""
	switch uniacid {
	case 21:
		supplyTokenUrl = fmt.Sprintf(conf.StatisC.PHPApi.SupplyAppletsAccessToken, constant.RoamUniacidStr)
	case 46:
		supplyTokenUrl = fmt.Sprintf(conf.StatisC.PHPApi.SupplyAppletsAccessToken, constant.ShareUniacidStr)
	default:
		supplyTokenUrl = fmt.Sprintf(conf.StatisC.PHPApi.SupplyAppletsAccessToken, constant.RoamUniacidStr)
	}
	signStr := util.Md5String(constant.SaburoKey + constant.SignKey)
	form := url.Values{}
	form.Add("saburoKey", constant.SaburoKey)
	form.Add("sign", signStr)
	phpRes, err := http.PostForm(supplyTokenUrl, form)
	if err != nil {
		return "", err
	}

	defer phpRes.Body.Close()

	phpAccessToken := struct {
		Code        int    `json:"code"`
		AccessToken string `json:"accessToken"`
	}{}

	if err := json.NewDecoder(phpRes.Body).Decode(&phpAccessToken); err != nil {
		return "", err
	}

	if phpAccessToken.AccessToken != "" {
		return phpAccessToken.AccessToken, nil
	}

	// next get access token form wxapi
	accessUrl := ""

	switch uniacid {
	case 21:
		accessUrl = fmt.Sprintf(conf.StatisC.WXApi.AppletsAccessTokenURL, conf.StatisC.WXSecret.RoamAppletsAppID, conf.StatisC.WXSecret.RoamAppletsAppSecret)
	case 46:
		accessUrl = fmt.Sprintf(conf.StatisC.WXApi.AppletsAccessTokenURL, conf.StatisC.WXSecret.ShareAppletsAppID, conf.StatisC.WXSecret.ShareAppletsAppSecret)
	}

	if accessUrl != "" {
		res, err := http.Get(accessUrl)

		defer res.Body.Close()

		if err != nil {
			return "", err
		}

		getAccessRes := struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
			ErrCode     int    `json:"errcode"`
			ErrMsg      string `json:"errmsg"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&getAccessRes); err != nil {
			return "", err
		}

		if getAccessRes.ErrCode != 0 {
			return "", fmt.Errorf(strconv.Itoa(getAccessRes.ErrCode) + "____" + getAccessRes.ErrMsg)
		}

		switch uniacid {
		case 21:
			err := components.Redis.Set(constant.RedisRoamAccessToken, getAccessRes.AccessToken, time.Duration(getAccessRes.ExpiresIn)*time.Second).Err()
			if err != nil {
				return getAccessRes.AccessToken, nil
			}
		case 46:
			err := components.Redis.Set(constant.RedisShareAccessToken, getAccessRes.AccessToken, time.Duration(getAccessRes.ExpiresIn)*time.Second).Err()
			if err != nil {
				return getAccessRes.AccessToken, nil
			}
		}

	}

	return "", nil

}

func getUsersSummary(token, beginDate, endDate string) (yNewUser, yCancelUser int, err error) {
	if token == "" || beginDate == "" || endDate == "" {
		return 0, 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}
	summaryUrl := fmt.Sprintf(conf.StatisC.WXApi.AppletsUsersSummary, token)
	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}
	byt, err := json.Marshal(req)
	if err != nil {
		return 0, 0, err
	}

	res, err := http.Post(summaryUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if err != nil {
		return 0, 0, err
	}

	defer res.Body.Close()

	getSummary := struct {
		ErrCode int                        `json:"errcode"`
		ErrMsg  string                     `json:"errmsg"`
		List    []protocol.SummaryListItem `json:"list"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&getSummary); err != nil {
		return 0, 0, err
	}

	if getSummary.ErrCode != 0 {
		return 0, 0, fmt.Errorf(strconv.Itoa(getSummary.ErrCode) + "____" + getSummary.ErrMsg)
	}
	yNewUser = 0
	yCancelUser = 0
	for _, summary := range getSummary.List {
		yNewUser += summary.NewUser
		yCancelUser += summary.CancelUser
	}
	return yNewUser, yCancelUser, nil
}

func getUsersCumulate(token, beginDate, endDate string) (yCumulateUsers int, err error) {

	if token == "" || beginDate == "" || endDate == "" {
		return 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}

	cumulateUrl := fmt.Sprintf(conf.StatisC.WXApi.AppletsUsersCumulate, token)
	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}

	byt, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	res, err := http.Post(cumulateUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	getCumulate := struct {
		ErrCode int                         `json:"errcode"`
		ErrMsg  string                      `json:"errmsg"`
		List    []protocol.CumulateListItem `json:"list"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&getCumulate); err != nil {
		return 0, err
	}

	if getCumulate.ErrCode != 0 {
		return 0, fmt.Errorf(strconv.Itoa(getCumulate.ErrCode) + "____" + getCumulate.ErrMsg)
	}

	return getCumulate.List[0].CumulateUser, nil

}
