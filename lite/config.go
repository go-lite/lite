package lite

type Config interface {
	Host() string
	Port() int
	Prefix() string
}
