package store

type Patch[T any] struct {
	Set   bool
	Value *T
}
