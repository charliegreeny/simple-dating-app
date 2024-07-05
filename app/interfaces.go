package app

import (
	"context"
)

type Creator[I, T any] interface {
	Create(ctx context.Context, input I) (T, error)
}

type IDGetter[T any] interface {
	Get(ctx context.Context, id string) (T, error)
}

type GetterCreator[I, T any] interface {
	Creator[I, T]
	IDGetter[T]
}

type Updater[I, T any] interface {
	Update(ctx context.Context, input I) (T, error)
}

type EntityService[I, T any] interface {
	Creator[I, T]
	IDGetter[T]
	Updater[I, T]
}

type Cache[K, V any] interface {
	Get(ctx context.Context, key K) (V, error)
	GetAll(ctx context.Context) []V
	Add(ctx context.Context, key K, v V) error
}
