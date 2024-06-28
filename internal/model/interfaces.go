package model

type Creator[I, T any] interface {
	Create(input I) (T, error)
}
