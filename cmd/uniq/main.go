package main

import (
	"log"
	"os"

	"hw1_tp/pkg/uniq"
)

func main() {
	uniqCommand, err := uniq.ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("error in main(): %v", err)
	}

	defer func() {
		err = uniqCommand.Close()
		if err != nil {
			log.Printf("error in main(): %v\n", err)
		}
	}()

	err = uniqCommand.Run()
	if err != nil {
		log.Printf("error in main(): %v", err)
	}
}
