package main

import (
	"fmt"
	"log"
	"os"

	"hw1_tp/pkg/calc"
)

func main() {
	defer func() {
		result, err := calc.Run(os.Args)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		log.Println("Result:", result)
	}()

}
