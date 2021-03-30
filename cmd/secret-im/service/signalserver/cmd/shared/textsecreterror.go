package shared
var (
	ErrAdxCheck=NewCodeError(300,"adx jwt check fail")
	ErrAdxHeadInvalid=NewCodeError(301,"miss head（name） or name value is empty")
)
