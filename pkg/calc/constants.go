package calc

const (
	operations = "+-/*("
	separators = operations + ")"

	// use -1 because ")" handle separately
	countOperations = len(operations)

	countOperandsForOperator = 2
)
