package main

import (
	"hw1_tp/pkg/calc"
	"log"
	"os"

	_ "hw1_tp/pkg/calc"
)

const correctCountArgs = 2

func main() {
	if len(os.Args) != correctCountArgs {
		log.Fatalf("unexpected count arguments %v", len(os.Args))
	}

	input := os.Args[1]

	result, err := calc.Calc(input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperations, calc.GetMapPriority())
	if err != nil {
		log.Fatalf("in main(): Error is: %v", err)
	}

	log.Println(result)
}
