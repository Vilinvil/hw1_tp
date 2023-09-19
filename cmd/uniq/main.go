package main

import (
	"fmt"
	"log"

	"hw1_tp/internal/uniq"
)

func main() {
	uniqCommand, err := uniq.ParseArgs()
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
