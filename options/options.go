package options

const (
	Internal = iota
	User     = iota
	No       = iota
)

type Options struct {
	Opt int
}
