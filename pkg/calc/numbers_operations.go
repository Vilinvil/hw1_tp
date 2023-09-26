package calc

import (
	"strconv"
)

type IsNumber func(token string) bool

type IsOperator func(token string) bool

type MapPriority map[string]int8

func IsInt(token string) bool {
	_, err := strconv.Atoi(token)

	return err == nil
}

func IsBasicOperator(token string) bool {
	for _, operation := range operators {
		if string(operation) == token {
			return true
		}
	}

	return false
}

func GetMapPriority() MapPriority {
	mapOperations := make(MapPriority, countOperations)
	mapOperations["("] = 1
	mapOperations["+"] = 2
	mapOperations["-"] = 2
	mapOperations["/"] = 3
	mapOperations["*"] = 3

	return mapOperations
}
