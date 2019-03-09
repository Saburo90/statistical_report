package server

import (
	"context"
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"gitee.com/NotOnlyBooks/statistical_report/models/users_model"
	"gitee.com/NotOnlyBooks/statistical_report/protos/statistical"
	"time"
)

type grpcServer struct {
}

func NewGRPCStatisticalServer() daily_statistical.StatisticalServer {
	return &grpcServer{}
}

func (s *grpcServer) GetUsersOverview(cotxt context.Context, req *daily_statistical.GetUsersOverviewRep) (*daily_statistical.GetUsersOverviewResp, error) {
	session := components.NewDBSession()

	defer session.Close()

	eweiShopMemMode := users_model.NewEweiShopMemberModel(session)

	nowTimeStr := time.Now().Format("2006-01-02")
	nt, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	yesterdayStart := nt.AddDate(0, 0, -1).Unix()
	yesterdayEnd := nt.Unix() - 1

	tPUsers, err := eweiShopMemMode.GetPlatRUsersNum(constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		return nil, err
	}

	tRusers, err := eweiShopMemMode.GetRoamAppletsUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		return nil, err
	}

	tRMusers, err := eweiShopMemMode.GetRoamMiniUsersNum("", "", constant.RoamingApplets, yesterdayEnd)

	if err != nil {
		return nil, err
	}

	tRNusers, tPNusers, tSNusers, err := eweiShopMemMode.GetYesterdayPNewUsersNum(constant.RoamingApplets, constant.ShareApplets, "", yesterdayStart, yesterdayEnd)

	if err != nil {
		return nil, err
	}

	tRANusers, tRMNusers, err := eweiShopMemMode.GetRapltAndRminiNewUsersNum(constant.RoamingApplets, "", "", yesterdayStart, yesterdayEnd)

	if err != nil {
		return nil, err
	}

	return &daily_statistical.GetUsersOverviewResp{
		TotalPlatformUsers: tPUsers,
		TotalRoamUsers:     tRusers,
		TotalRoamMiniUsers: tRMusers,
		TotalPlatnewUsers:  tPNusers,
		TotalRoamnewUsers:  tRNusers,
		TotalSharenewUsers: tSNusers,
		TotalRapltnewUsers: tRANusers,
		TotalRmininewUsers: tRMNusers,
	}, nil
}

func (s *grpcServer) GetWxAPIData(contxt context.Context, req *daily_statistical.GetWxAPIDataReq) (*daily_statistical.GetWxAPIDataResp, error) {

	return nil, nil
}

func (s *grpcServer) GetSalesData(contxt context.Context, req *daily_statistical.GetSalesDataReq) (*daily_statistical.GetSalesDataResp, error) {
	return nil, nil
}

func (s *grpcServer) GetRoamMiniWxAPIData(contxt context.Context, req *daily_statistical.GetRoamMiniWxAPIDataReq) (*daily_statistical.GetRoamMiniWxAPIDataResp, error) {
	return nil, nil
}

func appletsWxAPISummary(begineDate, endDate, token string) (int64, error) {
	//url := fmt.Sprintf(conf.StatisC.WXApi.AppletsUsersCumulate, token)

	return 0, nil
}
