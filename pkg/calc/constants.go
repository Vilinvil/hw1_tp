package calc

const (
	startInsteadUnaryMinus = "(0"

	correctCountArgs = 2

	operators  = "+-*/"
	operations = operators + "("
	separators = operations + ")"

	countOperations = len(operations)

	countOperandsForOperator = 2
)

func getSlFromStr(from string) []string {
	result := make([]string, len(from))

	for i, v := range from {
		result[i] = string(v)
	}

	return result
}
