package tokens

type Claims interface {
	Valid() bool
}
