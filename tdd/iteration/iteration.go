package iteration

import "strings"

// Repeat return a (character) repited many (times)
func Repeat(character string, times int) string {
	var str strings.Builder

	for i := 0; i < times; i++ {
		str.WriteString(character)
		//repeated = strings.Join([]string{repeated, character}, "")
	}

	return str.String()
}
