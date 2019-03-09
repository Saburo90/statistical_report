package wxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"gitee.com/NotOnlyBooks/statistical_report/protocol"
	"net/http"
	"strconv"
)

// 调用微信获取用户基本信息API获取用户unionid(单条版)
func GetAppletUnionidByWxAPI(apltToken, apltOpenid string) (string, error) {
	apltUnionidUrl := fmt.Sprintf(conf.StatisC.WXApi.AppletsUserInfo, apltToken, apltOpenid)

	res, err := http.Get(apltUnionidUrl)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	getUserUnionidRes := struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Subscribe int    `json:"subscribe"`
		Openid    string `json:"openid"`
		Unionid   string `json:"unionid"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&getUserUnionidRes); err != nil {
		return "", err
	}

	if getUserUnionidRes.ErrCode != 0 {
		return "", fmt.Errorf(strconv.Itoa(getUserUnionidRes.ErrCode) + "----" + getUserUnionidRes.ErrMsg)
	}

	if getUserUnionidRes.Subscribe == 0 {
		return "", fmt.Errorf(strconv.Itoa(getUserUnionidRes.ErrCode) + "---- user not subscribed")
	}
	return getUserUnionidRes.Unionid, nil
}

// 调用微信批量获取用户基本信息API获取多用户unionid(批量版)
func GetAppletUnionidsByWxAPI(openidList protocol.GetUsersInfos, apltToken string) (protocol.UsersInfosList, error) {
	var unionidList protocol.UsersInfosList

	batchgetUrl := fmt.Sprintf(conf.StatisC.WXApi.AppletsUsersInfos, apltToken)

	byt, err := json.Marshal(openidList)

	if err != nil {
		return unionidList, err
	}

	res, err := http.Post(batchgetUrl, "application/json; charset=utf-8", bytes.NewReader(byt))

	if err != nil {
		return unionidList, err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&unionidList); err != nil {
		return unionidList, err
	}

	if unionidList.ErrCode != 0 {
		return unionidList, fmt.Errorf(strconv.Itoa(unionidList.ErrCode) + "----" + unionidList.ErrMsg)
	}

	return unionidList, nil
}
