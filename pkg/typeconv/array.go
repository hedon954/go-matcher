package typeconv

func MapToSlice[T comparable, V any](m map[T]V) []T {
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func SliceToMap[T comparable](s []T) map[T]bool {
	res := make(map[T]bool, len(s))
	for _, v := range s {
		res[v] = true
	}
	return res
}
