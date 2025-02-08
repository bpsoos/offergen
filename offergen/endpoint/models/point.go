package models

func Point[T any](v T) *T {
	return &v
}
