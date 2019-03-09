package timed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/components/grpccli"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/models/statis_model"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"gitee.com/NotOnlyBooks/statistical_report/timed/mini_token_pb"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultExpire = time.Duration(24 * 60 * 60)
)

func GetMiniWxAPIData() {
	var (
		miniToken    string
		yVisitTotal  int
		visitPv      int
		visitUv      int
		visitUvNew   int
		stayTimeUv   float64
		rpageVisitPv int
		rpageVisitUv int
		dpageVisitPv int
		dpageVisitUv int
		dailyCount   int64
		searchPv     int64
		searchUv     int64
		meanSear     decimal.NullDecimal
		err          error
	)

	zap.L().Info("开始获取漫游鲸小程序access_token")
	for retry := 5; retry > 0; retry-- {
		miniToken, err = getAccessToken()

		if err != nil {
			zap.L().Error("获取漫游鲸小程序access_token失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸小程序access_token失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			return
		}
	}

	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("20060102")

	zap.L().Info("开始获取漫游鲸小程序截至昨日终止时间总访问人数")
	for retry := 3; retry > 0; retry-- {
		yVisitTotal, err = getAnalyDailySummary(miniToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游鲸小程序截至昨日终止时间访问人数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸小程序截至昨日终止时间访问人数失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			return
		}
	}

	zap.L().Info("开始获取漫游鲸小程序昨日访问数据")
	for retry := 3; retry > 0; retry-- {
		visitPv, visitUv, visitUvNew, stayTimeUv, err = getAnalyDailyVisitTrend(miniToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游鲸小程序昨日访问数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸小程序昨日访问数据失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			return
		}
	}

	zap.L().Info("开始获取漫游鲸小程序(回收/商品详情)页访问数据")
	for retry := 3; retry > 0; retry-- {
		rpageVisitPv, rpageVisitUv, dpageVisitPv, dpageVisitUv, err = getAnalyVisitPage(miniToken, yesterdayStr, yesterdayStr)

		if err != nil {
			zap.L().Error("获取漫游鲸小程序(回收/商品详情)页访问数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(time.Second * 3)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取漫游鲸小程序(回收/商品详情)页访问数据失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			return
		}

	}

	session := components.NewDBSession()

	defer session.Close()

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1

	pvSearch := statis_model.SearchGoods{}

	zap.L().Info("开始获取昨日搜索商品次数")
	for retry := 3; retry > 0; retry-- {
		searchPv, err = session.Where("search_time > ?", yesterdayStart).
			And("search_time <= ?", yesterdayEnd).
			And("uniacid = ?", constant.RoamingApplets).
			Count(&pvSearch)

		if err != nil {
			zap.L().Error("获取昨日搜索商品次数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取昨日搜索商品次数失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			break
		}

	}

	zap.L().Info("开始获取昨日搜索商品用户数")
	uvSearch := statis_model.SearchGoods{}
	for retry := 3; retry > 0; retry-- {
		searchUv, err = session.Distinct("uid").
			Where("search_time > ?", yesterdayStart).
			And("search_time <= ?", yesterdayEnd).
			And("uniacid = ?", constant.RoamingApplets).
			Count(&uvSearch)

		if err != nil {
			zap.L().Error("获取昨日搜索商品用户数失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
			time.Sleep(3 * time.Second)
		} else {
			break
		}

		if retry == 0 {
			zap.L().Error("获取昨日搜索商品用户数失败，重试机会已耗尽，今日获取漫游鲸小程序数据分析失败")
			break
		}

	}

	if searchUv > 0 {
		zap.L().Info("开始计算昨日人均搜索商品次数")
		meanSear.Decimal = decimal.New(searchPv, 0).Div(decimal.New(searchUv, 0)).Round(2)
		meanSear.Valid = true
	}

	staticticalModel := statis_model.NewStatisticalReportModel(session)

	nyesStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	zap.L().Info("开始检验是否存在昨日统计记录")
	for retry := 3; retry > 0; retry-- {

		dailyCount, err = staticticalModel.CheckDaily(nyesStr)

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
		zap.L().Info("开始记录漫游鲸小程序访问统计数据")
		miniStatis := statis_model.MiniStatis{
			Daily:                  nyesStr,
			TotalRminiVisitsU:      int64(yVisitTotal),
			TotalRminiNewVisitsUv:  int64(visitUv),
			TotalRminiNewVisitsNo:  int64(visitPv),
			TotalRminiNewVisitsUsr: int64(visitUvNew),
			TotalRminiNewRcyPageU:  int64(rpageVisitUv),
			TotalRminiNewRcyPageN:  int64(rpageVisitPv),
			TotalRminiNewBdePageU:  int64(dpageVisitUv),
			TotalRminiNewBdePageN:  int64(dpageVisitPv),
			TotalRminiNewVitStay:   stayTimeUv,
			TotalRminiNewSearchPv:  searchPv,
			TotalRminiNewSearchUv:  searchUv,
			Meansearch:             meanSear,
		}
		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				InsertOne(&miniStatis)

			if err != nil {
				fmt.Println("err = ", err)
				zap.L().Error("记录漫游鲸小程序访问数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("漫游鲸小程序访问数据记录成功", zap.String("记录日历", nyesStr), zap.String("记录时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("记录漫游鲸小程序访问数据失败，重试机会已耗尽，今日漫游鲸小程序访问数据获取失败")
				return
			}
		}

	} else {
		// 更新记录
		zap.L().Info("开始更新漫游鲸小程序访问数据")
		miniStatisUpd := statis_model.MiniStatis{
			Daily:                  nyesStr,
			TotalRminiVisitsU:      int64(yVisitTotal),
			TotalRminiNewVisitsUv:  int64(visitUv),
			TotalRminiNewVisitsNo:  int64(visitPv),
			TotalRminiNewVisitsUsr: int64(visitUvNew),
			TotalRminiNewRcyPageU:  int64(rpageVisitUv),
			TotalRminiNewRcyPageN:  int64(rpageVisitPv),
			TotalRminiNewBdePageU:  int64(dpageVisitUv),
			TotalRminiNewBdePageN:  int64(dpageVisitPv),
			TotalRminiNewVitStay:   stayTimeUv,
			TotalRminiNewSearchPv:  searchPv,
			TotalRminiNewSearchUv:  searchUv,
			Meansearch:             meanSear,
		}

		for retry := 5; retry > 0; retry-- {
			_, err := session.Table("ims_statistical_report").
				Where("daily = ?", nyesStr).
				Update(&miniStatisUpd)

			if err != nil {
				zap.L().Error("更新漫游鲸小程序访问数据失败, 3秒后重试", zap.Error(err), zap.Int("重试机会", retry))
				time.Sleep(3 * time.Second)
			} else {
				zap.L().Info("漫游鲸小程序访问数据更新成功", zap.String("更新日历", nyesStr), zap.String("更新时间", time.Now().Format("20060102150405")))
				return
			}

			if retry == 0 {
				zap.L().Error("更新漫游鲸小程序访问数据失败，重试机会已耗尽，今日漫游鲸小程序访问数据获取失败")
				return
			}
		}
	}
}

func getAccessToken() (string, error) {
	mitokenConn, err := grpccli.GetgRPCCli(constant.GRpcMiniToken)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	accTokenCli := pb_brilliant.NewAccessTokenClient(mitokenConn)

	accData, err := accTokenCli.GetAccessToken(context.Background(), &pb_brilliant.GetAccessTokenReq{Name: "roam_whale_mini"})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return accData.Token, nil
}

func getAnalyDailySummary(token, beginDate, endDate string) (int, error) {

	if token == "" || beginDate == "" || endDate == "" {
		return 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}
	dailySummaryUrl := fmt.Sprintf(conf.StatisC.WXApi.MiniAnaDailySummary, token)

	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}
	byt, err := json.Marshal(req)

	if err != nil {
		return 0, err
	}

	res, err := http.Post(dailySummaryUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	getDailySummary := struct {
		ErrCode int                         `json:"errcode"`
		ErrMsg  string                      `json:"errmsg"`
		List    []protocol.DailySummaryItem `json:"list"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&getDailySummary); err != nil {
		return 0, err
	}

	fmt.Println("getDailySummary = ", getDailySummary.List)

	if getDailySummary.ErrCode != 0 {
		return 0, fmt.Errorf(strconv.Itoa(getDailySummary.ErrCode) + "____" + getDailySummary.ErrMsg)
	}

	if len(getDailySummary.List) == 0 {
		return 0, fmt.Errorf("There Is No Daily Summary Data")
	}
	return getDailySummary.List[0].VisitTotal, nil
}

func getAnalyDailyVisitTrend(token, beginDate, endDate string) (visit_pv, visit_uv, visit_uv_new int, stay_time_uv float64, err error) {
	if token == "" || beginDate == "" || endDate == "" {
		return 0, 0, 0, 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}
	dailyVisitTrendUrl := fmt.Sprintf(conf.StatisC.WXApi.MiniAnaVisitTrend, token)
	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}
	byt, errs := json.Marshal(req)

	if errs != nil {
		return 0, 0, 0, 0, errs
	}

	res, errs := http.Post(dailyVisitTrendUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if errs != nil {
		fmt.Println("visitErrs = ", err)
		return 0, 0, 0, 0, errs
	}

	defer res.Body.Close()

	getDailyVisitTrend := struct {
		ErrCode int                            `json:"errcode"`
		ErrMsg  string                         `json:"errmsg"`
		List    []protocol.DailyVisitTrendItem `json:"list"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&getDailyVisitTrend); err != nil {
		fmt.Println("visitErr = ", err)
		return 0, 0, 0, 0, err
	}

	if getDailyVisitTrend.ErrCode != 0 {
		return 0, 0, 0, 0, fmt.Errorf(strconv.Itoa(getDailyVisitTrend.ErrCode) + "____" + getDailyVisitTrend.ErrMsg)
	}

	if len(getDailyVisitTrend.List) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("There Is No Daily Visit Data")
	}

	components.Redis.Set(constant.RedisYesNewVisitUv, getDailyVisitTrend.List[0].VisitUv, defaultExpire*time.Second)
	components.Redis.Set(constant.RedisYesNewVisitUvNew, getDailyVisitTrend.List[0].VisitUvNew, defaultExpire*time.Second)

	return getDailyVisitTrend.List[0].VisitPv, getDailyVisitTrend.List[0].VisitUv, getDailyVisitTrend.List[0].VisitUvNew, getDailyVisitTrend.List[0].StayTimeUv, nil
}

func getAnalyVisitPage(token, beginDate, endDate string) (rpage_visit_pv, rpage_visit_uv, dpage_visit_pv, dpage_visit_uv int, err error) {
	if token == "" || beginDate == "" || endDate == "" {
		return 0, 0, 0, 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}
	visitPageUrl := fmt.Sprintf(conf.StatisC.WXApi.MiniAnaVisitPage, token)
	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}
	byt, errs := json.Marshal(req)

	if errs != nil {
		return 0, 0, 0, 0, errs
	}

	res, errs := http.Post(visitPageUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if errs != nil {
		return 0, 0, 0, 0, errs
	}

	defer res.Body.Close()

	getVisitPage := struct {
		ErrCode int                      `json:"errcode"`
		ErrMsg  string                   `json:"errmsg"`
		RefDate string                   `json:"ref_date"`
		List    []protocol.VisitPageItem `json:"list"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&getVisitPage); err != nil {
		return 0, 0, 0, 0, err
	}

	if getVisitPage.ErrCode != 0 {
		return 0, 0, 0, 0, fmt.Errorf(strconv.Itoa(getVisitPage.ErrCode) + "____" + getVisitPage.ErrMsg)
	}

	if len(getVisitPage.List) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("There Is No Visit Page Data")
	}

	for _, item := range getVisitPage.List {
		if item.PagePath == "pages/secondHandBook/index" {
			rpage_visit_pv = item.PageVisitPv
			rpage_visit_uv = item.PageVisitUv
		}

		if item.PagePath == "pages/home/bookInfo/index" {
			dpage_visit_pv = item.PageVisitPv
			dpage_visit_uv = item.PageVisitUv
		}
	}

	components.Redis.Set(constant.RedisRcyPageVisitUv, rpage_visit_uv, defaultExpire*time.Second)
	// insert to redis

	return rpage_visit_pv, rpage_visit_uv, dpage_visit_pv, dpage_visit_uv, nil
}
