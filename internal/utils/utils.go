package utils

// length of multidimensional array
func LMDA[V any](arr [][]V) int {
	res := 0
	for _, v := range arr {
		res += len(v)
	}
	return res
}

func IsChannelClosed(ch <-chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}
