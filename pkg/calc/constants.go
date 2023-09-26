package calc

const (
	operators  = "+-*/"
	operations = operators + "("
	separators = operations + ")"

	countOperations = len(operations)

	countOperandsForOperator = 2
)
