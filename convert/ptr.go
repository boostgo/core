package convert

func Ptr[T any](x T) *T {
	return &x
}
