package iteration

import (
	"strings"
)

func Repeat(a string, repeat int) string {
	var repeated strings.Builder
	for i := 0; i < repeat; i++ {
		repeated.WriteString(a)	
	}
	return repeated.String()
}