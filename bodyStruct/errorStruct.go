package bodyStruct

type ErrorStruct struct {
	Code       string `json:"code"`
	Path       string `json:"path"`
	RequestId  string `json:"requestId"`
	Error      string `json:"error"`
	Message    string `json:"devMessage"`
	SourceFile string `json:"sourceFile,omitempty"`
	// TODO 시간정보?
}
