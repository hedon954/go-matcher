package collection

// InSlice checks if target is in slice
func InSlice[T comparable](target T, slice []T) bool {
	for _, t := range slice {
		if target == t {
			return true
		}
	}
	return false
}
