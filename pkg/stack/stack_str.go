package stack

type StackStr []string

func (s *StackStr) Push(str string) {
	*s = append(*s, str)
}

func (s *StackStr) Pop() {
	*s = (*s)[:len(*s)-1]
}
