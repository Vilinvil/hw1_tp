package uniq

import "fmt"

func Run(args []string) error {
	uniqCommand, err := ParseArgs(args)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	defer func() {
		err = uniqCommand.Close()
		if err != nil {
			fmt.Printf("in defer in Run():\n %v", err)
		}
	}()

	err = uniqCommand.Run()
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}

	return nil
}
