package main

import (
	"fmt"
	"os"

	"hw1_tp/pkg/uniq"
)

func main() {
	defer func() {
		err := uniq.Run(os.Args[1:])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()
}
