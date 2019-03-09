package conf

type (
	WxAPI struct {
		AppletsAccessTokenURL string `toml:"get_applets_access_token"`
		AppletsUsersSummary   string `toml:"get_applets_users_summary"`
		AppletsUsersCumulate  string `toml:"get_applets_users_cumulate"`
		MiniAnaVisitPage      string `toml:"get_analysis_visit_page"`
		MiniAnaDailySummary   string `toml:"get_analysis_daliy_summary"`
		MiniAnaVisitTrend     string `toml:"get_analysis_visit_trend"`
		AppletsUserInfo       string `toml:"get_applets_user_info"`
		AppletsUsersInfos     string `toml:"get_applets_users_infos"`
	}

	WxSecret struct {
		RoamAppletsAppID      string `toml:"roam_applets_appId"`
		RoamAppletsAppSecret  string `toml:"roam_applets_appSecret"`
		ShareAppletsAppID     string `toml:"share_applets_appId"`
		ShareAppletsAppSecret string `toml:"share_applets_appSecret"`
	}
)
