package shared
var (
	ErrAdxCheck=NewCodeError(300,"adx jwt check fail")
	ErrAdxHeadInvalid=NewCodeError(301,"miss head（name） or name value is empty")

	ERRCODE_JSONMARSHAL = 302
	ERRCODE_JSONUNMARSHAL=303
	ERRCODE_SQLQUERY=304
	ERRCODE_SQLINSERT=305
	ERRCODE_WSCONNECTDUP=306
	ERRCODE_WSCONNECTERR=307

)
