package options

var BendSep = "$303"

const (
	Internal      = iota
	User          = iota
	Transpile     = iota
	Primitive     = iota
	ExecPrimitive = iota
	ExecBasm      = iota
	ShowTree      = true
	ExecBenv      = true
	InterpBenv    = false
	HideTree      = false
)
