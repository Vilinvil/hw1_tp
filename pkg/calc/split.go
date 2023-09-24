package calc

import (
	"bufio"
	"bytes"
	"strings"
)

func splitTokensInt(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	idxSeparator := bytes.IndexAny(data, separators)
	if idxSeparator == 0 {
		return 1, data[:1], nil
	} else if idxSeparator > 0 {
		return idxSeparator, data[:idxSeparator], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}

func splitIntoTokens(input string, splitFunc bufio.SplitFunc) []string {
	scanner := bufio.NewScanner(bytes.NewBufferString(input))

	scanner.Split(splitFunc)

	var result []string

	for scanner.Scan() {
		token := strings.TrimSpace(scanner.Text())
		if token != "" {
			result = append(result, token)
		}
	}

	return result
}
