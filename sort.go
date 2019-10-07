package cli

import "unicode"

// taken from https://github.com/urfave/cli/blob/master/sort_test.go
func lexicographicLess(i, j string) bool {
	iRunes := []rune(i)
	jRunes := []rune(j)

	lenShared := len(iRunes)
	if lenShared > len(jRunes) {
		lenShared = len(jRunes)
	}

	for index := 0; index < lenShared; index++ {
		ir := iRunes[index]
		jr := jRunes[index]

		if lir, ljr := unicode.ToLower(ir), unicode.ToLower(jr); lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return i < j
}
