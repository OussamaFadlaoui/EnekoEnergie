package helpers

func Contains(haystack []int, needle int) bool {
	for i := range haystack {
		if i == needle {
			return true
		}
	}

	return false
}

func ContainsAll(haystack []int, needles ... int) bool {
	flagAll := true

	for needle := range needles {
		if !Contains(haystack, needle) {
			flagAll = false
		}
	}

	return flagAll
}
