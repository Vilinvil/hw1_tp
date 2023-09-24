package main

import (
	"fmt"
	"hw1_tp/pkg/calc"

	_ "hw1_tp/pkg/calc"
)

func main() {
	input := "233 * 344 / (5 * 666) - 10 "

	fmt.Println(calc.Calc(input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperations, calc.GetMapPriority()))
}
