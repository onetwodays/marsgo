package shared

const (
	DEFAULTCODE = 1001
	// 注册的错误码
	DuplicateUsername =1002
	DuplicateMobile = 1003

	//登录的错误码
	UsernameUnRegister = 1004
	IncorrectPassword  = 1005

	//JWT校验
	GetUserInfoFailed =1006
	AuthDeny = 1007

	//获取用户信息
	UserNotFound =1008

	//参数错误
	ParamError =1009

)

var CodeErrorMap =map[int]string{
	DuplicateUsername:"用户名已经注册",
	DuplicateMobile:"手机号已经被占用",

	UsernameUnRegister:"用户未注册",
	IncorrectPassword: "用户密码错误",

	GetUserInfoFailed:"用户信息获取失败",
	AuthDeny:"用户信息不一致",

	UserNotFound: "用户不存在",
	ParamError: "参数错误",
}
var ErrorParam = NewCodeError(ParamError,CodeErrorMap[ParamError])





type CodeError struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int,msg string) error{
	return &CodeError{Code:code,Msg: msg}
}
func NewDefaultCodeError(msg string) error{
	return &CodeError{Code:DEFAULTCODE,Msg: msg}
}
func(e *CodeError) Error() string{
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse{
	return &CodeErrorResponse{
		Code:e.Code,
		Msg: e.Msg,
	}
}