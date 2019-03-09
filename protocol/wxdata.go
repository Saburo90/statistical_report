package protocol

type (
	PublicReq struct {
		BeginDate string `json:"begin_date"`
		EndDate   string `json:"end_date"`
	}

	SummaryListItem struct {
		RefDate    string `json:"ref_date"`
		UserSource int    `json:"user_source"`
		NewUser    int    `json:"new_user"`
		CancelUser int    `json:"cancel_user"`
	}

	CumulateListItem struct {
		RefDate      string `json:"ref_date"`
		CumulateUser int    `json:"cumulate_user"`
	}

	DailySummaryItem struct {
		RefDate    string `json:"ref_date"`
		VisitTotal int    `json:"visit_total"`
		SharePv    int    `json:"share_pv"`
		ShareUv    int    `json:"share_uv"`
	}

	DailyVisitTrendItem struct {
		RefDate         string  `json:"ref_date"`
		SessionCnt      int     `json:"session_cnt"`
		VisitPv         int     `json:"visit_pv"`
		VisitUv         int     `json:"visit_uv"`
		VisitUvNew      int     `json:"visit_uv_new"`
		StayTimeUv      float64 `json:"stay_time_uv"`
		StayTimeSession float64 `json:"stay_time_session"`
		VisitDepth      float64 `json:"visit_depth"`
	}

	VisitPageItem struct {
		PagePath       string  `json:"page_path"`
		PageVisitPv    int     `json:"page_visit_pv"`
		PageVisitUv    int     `json:"page_visit_uv"`
		PageStayTimePv float64 `json:"page_staytime_pv"`
		EntryPagePv    int     `json:"entrypage_pv"`
		ExitPagePv     int     `json:"exitpage_pv"`
		PageSharePv    int     `json:"page_share_pv"`
		PageShareUv    int     `json:"page_share_uv"`
	}

	GetWxDataResp struct {
		YRtotalUsers  int `json:"yesterday_roamtotal_fans"`
		YRnewUsers    int `json:"yesterday_roamnew_fans"`
		YRcancelUsers int `json:"yesterday_roamcancel_fans"`
		YStotalUsers  int `json:"yesterday_sharetotal_fans"`
		YSnewUsers    int `json:"yesterday_sharenew_fans"`
		YScancelUsers int `json:"yesterday_sharecancel_fans"`
	}

	GetMiniDataResp struct {
		MiniVisitTotal     int     `json:"visit_total"`
		MiniVisitPv        int     `json:"yesterday_visit_pv"`
		MiniVisitUv        int     `json:"yesterday_visit_uv"`
		MiniVisitUvNew     int     `json:"yesterday_visit_uvnew"`
		MiniStayTimeUv     float64 `json:"yesterday_stay_timeuv"`
		MiniRcyPageVisitPv int     `json:"rcy_page_visitpv"`
		MiniRcyPageVisitUv int     `json:"rcy_page_visituv"`
		MiniDetPageVisitPv int     `json:"det_page_visitpv"`
		MiniDetPageVisitUv int     `json:"det_page_visituv"`
	}

	OpenidList struct {
		Openid string `json:"openid"`
		Lang   string `json:"lang"`
	}

	UserInfoItem struct {
		Subscribe int    `json:"subscribe"`
		Openid    string `json:"openid"`
		Unionid   string `json:"unionid"`
	}

	UsersInfosList struct {
		ErrCode      int            `json:"errcode"`
		ErrMsg       string         `json:"errmsg"`
		UserInfoList []UserInfoItem `json:"user_info_list"`
	}

	GetUsersInfos struct {
		UserList []OpenidList `json:"user_list"`
	}
)
