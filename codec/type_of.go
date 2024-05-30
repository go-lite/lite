package codec

import "reflect"

type TypeOf interface {
	TypeOf() reflect.Type
}
