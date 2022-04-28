package options

var BendSep = "$303"

const (
	Internal        = iota
	UserTranslate   = iota
	Transpile       = iota
	Primitive       = iota
	InterpPrimitive = iota
	ExecBasm        = iota
	ShowTree        = true
	ExecBenv        = true
	InterpBenv      = false
	HideTree        = false
)
