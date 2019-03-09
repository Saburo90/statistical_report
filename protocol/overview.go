package protocol

type (
	OverviewReq struct {
		Operator    string `json:"operator"`
		OperateTime int64  `json:"operateTime"`
		ClientIP    string `json:"clientIP"`
		Sign        string `json:"sign"`
	}

	OverviewResp struct {
		TotalPUsers   int64 `json:"total_platform_users"`
		TotalRUsers   int64 `json:"total_roam_users"`
		TotalRMUsers  int64 `json:"total_roamMini_users"`
		TotalPNUsers  int64 `json:"total_platnew_users"`
		TotalRNUsers  int64 `json:"total_roamnew_users"`
		TotalSNUsers  int64 `json:"total_sharenew_users"`
		TotalRANUsers int64 `json:"total_rapltnew_users"`
		TotalRMNUsers int64 `json:"total_rmininew_users"`
	}
)
