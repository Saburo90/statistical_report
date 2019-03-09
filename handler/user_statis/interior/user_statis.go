package interior

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"gitee.com/NotOnlyBooks/statistical_report/util"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetAppletsAccessToken(uniacid uint64) (string, error) {
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

func GetUsersSummary(token, beginDate, endDate string) (yNewUser, yCancelUser int, err error) {
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

func GetUsersCumulate(token, beginDate, endDate string) (yCumulateUsers int, err error) {

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
