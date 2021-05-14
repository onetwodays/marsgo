package shared

type OkResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}
func NewOkResponse(data interface{}) *OkResponse {
	return &OkResponse{
		Code:0,
		Msg: "successful",
		Data: data,
	}
}

func NewErrResponse(code int,err error) *ErrResponse{
	return &ErrResponse{
		Code:code,
		Msg: err.Error(),
	}
}

func (this *ErrResponse ) Error() string{
	return ""
}