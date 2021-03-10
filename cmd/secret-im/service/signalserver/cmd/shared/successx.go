package shared

type OkResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
func NewOkResponse(data interface{}) *OkResponse {
	return &OkResponse{
		Code:0,
		Msg: "successful",
		Data: data,
	}
}