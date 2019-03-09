package interior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/components/grpccli"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"gitee.com/NotOnlyBooks/statistical_report/timed/mini_token_pb"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultExpire = time.Duration(24 * 60 * 60)
)

func GetAccessToken() (string, error) {
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

func GetAnalyDailySummary(token, beginDate, endDate string) (int, error) {

	fmt.Println("token = ", token)

	if token == "" || beginDate == "" || endDate == "" {
		return 0, fmt.Errorf("PARAMETER MUST NOT BE EMPTY")
	}
	dailySummaryUrl := fmt.Sprintf(conf.StatisC.WXApi.MiniAnaDailySummary, token)

	fmt.Println("url = ", dailySummaryUrl)

	req := protocol.PublicReq{
		BeginDate: beginDate,
		EndDate:   endDate,
	}
	byt, err := json.Marshal(req)

	if err != nil {
		return 0, err
	}

	res, err := http.Post(dailySummaryUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	fmt.Println("err = ", err)

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

func GetAnalyDailyVisitTrend(token, beginDate, endDate string) (visit_pv, visit_uv, visit_uv_new int, stay_time_uv float64, err error) {
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
	fmt.Println("vtrendUrl = ", dailyVisitTrendUrl)
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

func GetAnalyVisitPage(token, beginDate, endDate string) (rpage_visit_pv, rpage_visit_uv, dpage_visit_pv, dpage_visit_uv int, err error) {
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

	// insert to redis

	return rpage_visit_pv, rpage_visit_uv, dpage_visit_pv, dpage_visit_uv, nil
}
