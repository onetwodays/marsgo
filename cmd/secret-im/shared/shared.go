package shared

const (
	DEFAULTCODE = 1001
	DuplicateUsername =1002
	DuplicateMobile = 1003
)

var CodeErrorMap =map[int]string{
	DuplicateUsername:"用户名已经注册",
	DuplicateMobile:"手机号已经被占用",
}





type codeError struct {
	code int `json:"code"`
	msg string `json:"msg"`
}
func NewCodeError(code int,msg string) error{
	return &codeError{code:code,msg: msg}
}
func NewDefaultCodeError(msg string) error{
	return &codeError{code:DEFAULTCODE,msg: msg}
}
func(e *codeError) Error() string{
	return e.msg
}