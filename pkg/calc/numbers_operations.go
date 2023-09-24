package calc

import (
	"strconv"
)

type isNumber func(token string) bool

type isOperations func(token string) bool

type mapPriority map[string]int8

func isInt(token string) bool {
	_, err := strconv.Atoi(token)

	return err == nil
}

func isBasicOperations(token string) bool {
	for _, operation := range operations {
		if string(operation) == token {
			return true
		}
	}

	return false
}

func getMapPriority() mapPriority {
	mapOperations := make(mapPriority, countOperations)
	mapOperations["("] = 1
	mapOperations["+"] = 2
	mapOperations["-"] = 2
	mapOperations["/"] = 3
	mapOperations["*"] = 3

	return mapOperations
}

type Number interface {
	FromString(string) (Number, error)
	Sum(another Number) Number
	Sub(another Number) Number
	Mul(another Number) Number
	Div(another Number) Number
}

type MyInt int

func (m MyInt) FromString(from string) (MyInt, error) {
	res, err := strconv.Atoi(from)

	return MyInt(res), err
}

func (m MyInt) Sum(another MyInt) MyInt {
	return m + another
}

func (m MyInt) Sub(another MyInt) MyInt {
	return m - another
}

func (m MyInt) Mul(another MyInt) MyInt {
	return m * another
}

func (m MyInt) Div(another MyInt) MyInt {
	return m / another
}
