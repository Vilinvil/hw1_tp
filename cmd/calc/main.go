package main

import (
	"fmt"

	_ "hw1_tp/pkg/calc"
	"hw1_tp/pkg/stack"
)

func main() {
	st := stack.Stack[string]{}
	fmt.Println(st)
	st.Push("asfd")
	st.Push("asfd")
	st.Push("asfd")
	st.Push("asfd")
	st.Push("asfd")
	fmt.Println(st)
}
