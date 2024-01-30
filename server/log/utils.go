package log

const (
	timeFormat        = "2006-01-02 15:04:05-0700"
	termTimeFormat    = "01-02|15:04:05.000"
	floatFormat       = 'f'
	termMsgJust       = 40
	termCtxMaxPadding = 40
	skipLevel         = 3
	LogLevel          = 4
)

const (
	LvlCrit int = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
	LvlTrace
)
