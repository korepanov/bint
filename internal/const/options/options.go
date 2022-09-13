package options

var BendSep = "$303"
var MaxTrashLen = 5

const (
	Internal         = iota
	UserTranslate    = iota
	Transpile        = iota
	Primitive        = iota
	InterpPrimitive  = iota
	ExecBasm         = iota
	Encrypt          = iota
	ExecEncrypt      = iota
	UserValidate     = iota
	InternalValidate = iota
	ShowTree         = true
	ExecBenv         = true
	InterpBenv       = false
	HideTree         = false
)
