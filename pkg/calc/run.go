package calc

import "fmt"

func Run(args []string) (float64, error) {
	if len(args) != correctCountArgs {
		return 0, fmt.Errorf(errTemplate, ErrUnexpectedCountArgs)
	}

	input := args[1]
	if input == "--help" || input == "-h" {
		return 0, fmt.Errorf(errTemplate, ErrHelpMessage)
	}

	result, err := Calc(input, SplitTokensInt, IsFloat, IsBasicOperator, GetMapPriority())
	if err != nil {
		return 0, fmt.Errorf(errTemplate, err)
	}

	return result, nil
}
