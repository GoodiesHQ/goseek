package utils

func Remove[T comparable](slice []T, index int) []T {
	if index < 0 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func Find[T comparable](slice []T, value T) int {
	for i, val := range slice {
		if val == value {
			return i
		}
	}
	return -1
}

func RemoveAll[T comparable](slice []T, values ...T) []T {
	for _, value := range values {
		slice = Remove(slice, Find(slice, value))
	}
	return slice
}