package codec

type ContentType interface {
	ContentType() string
	StructTag() string
}
