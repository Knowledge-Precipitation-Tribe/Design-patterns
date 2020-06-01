package apiRequest

type ApiRequest struct {
	baseURL string
	appID string
	timeStamp int
	token string
}

func GenFakeReq(url string) ApiRequest {
	return ApiRequest{}
}