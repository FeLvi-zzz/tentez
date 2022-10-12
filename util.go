package tentez

func chunk[T any](arr []T, size int) [][]T {
	res := [][]T{}
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		res = append(res, arr[i:end])
	}
	return res
}
