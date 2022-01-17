package bodyStruct

type ErrorStruct struct {
	Code      string `json:"code"`
	Path      string `json:"path"`
	RequestId string `json:"requestId"`
	Error     string `json:"error"`
	Message   string `json:"devMessage"`
	// TODO 시간정보?
}
