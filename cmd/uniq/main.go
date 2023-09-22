package main

import (
	"fmt"
	"log"
	"os"

	"hw1_tp/pkg/uniq"
)

func main() {
	uniqCommand, err := uniq.ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Errorf("error in main(): %w", err))
	}

	defer func() {
		err = uniqCommand.Close()
		if err != nil {
			log.Println(fmt.Errorf("error in main(): %w", err))
		}
	}()

	err = uniqCommand.Run()
	if err != nil {
		log.Println(fmt.Errorf("error in main(): %w", err))
	}
}
