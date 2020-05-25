package log

// Level of severity.严重性
type Level int

// Verbose(详情) is a boolean type that implements Info, Infov (like Printf) etc.
type Verbose bool

// common log level.
const (
	_debugLevel Level = iota
	_infoLevel
	_warnLevel
	_errorLevel
	_fatalLevel
)

var levelNames = [...]string{
	_debugLevel: "DEBUG",
	_infoLevel:  "INFO",
	_warnLevel:  "WARN",
	_errorLevel: "ERROR",
	_fatalLevel: "FATAL",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}
