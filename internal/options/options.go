package options

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
