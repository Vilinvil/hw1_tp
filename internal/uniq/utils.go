package uniq

import "strings"

func calcCountTrueFlags(flags ...bool) int {
	countFlags := 0

	for _, value := range flags {
		if value {
			countFlags++
		}
	}

	return countFlags
}

func TrimFirstFields(str string, count uint) string {
	idx := strings.Index(str, " ")
	for ; count > 0; count-- {
		if idx == -1 {
			return ""
		}

		str = str[idx+1:]
		idx = strings.Index(str, " ")
	}

	return str
}
