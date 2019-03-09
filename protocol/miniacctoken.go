package protocol

type (
	GetAccessTokenReq struct {
		Name string `json:"name"`
	}
	GetAccessTokenResp struct {
		Token string `json:"token"`
	}
)
