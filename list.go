package lite

type List[T any] struct {
	Items  []T `json:"items"`
	Length int `json:"length"`
}

func NewList[T any](items []T) List[T] {
	return List[T]{
		Items:  items,
		Length: len(items),
	}
}
