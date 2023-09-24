package stack

type Stack[T any] struct {
	buf []T
}

func (s *Stack[any]) Push(el any) {
	s.buf = append(s.buf, el)
}

func (s *Stack[any]) Pop() {
	s.buf = s.buf[:len(s.buf)-1]
}

func (s *Stack[any]) Top() any {
	return s.buf[len(s.buf)-1]
}

func (s *Stack[any]) Len() int {
	return len(s.buf)
}
