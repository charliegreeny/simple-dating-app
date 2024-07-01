package model

type Creator[I, T any] interface {
	Create(input I) (T, error)
}

type Getter[T any] interface {
	Get(id string) (T, error)
}

type GetterCreator[I, T any] interface {
	Creator[I, T]
	Getter[T]
}
