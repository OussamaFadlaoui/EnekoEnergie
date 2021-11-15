package helpers

func ArrContains(haystack []int, needle int) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}

	return false
}

func ArrContainsAll(haystack []int, needles ... int) bool {
	for _, needle := range needles {
		if !ArrContains(haystack, needle) {
			return false
		}
	}

	return true
}

func ArrContainsAny(haystack []int, needles ... int) bool {
	for _, needle := range needles {
		if ArrContains(haystack, needle) {
			return true
		}
	}

	return false
}

func ArrContainsNone(haystack []int, needles ... int) bool {
	return !ArrContainsAny(haystack, needles...)
}

func RemoveFromArr(haystack []int, needle int) []int {
	if needle < 0 || len(haystack) >= needle {
		return haystack
	}

	haystack[needle] = haystack[len(haystack)-1]
	return haystack[:len(haystack)-1]
}
