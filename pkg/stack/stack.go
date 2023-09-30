package stack

type Stack[T any] struct {
	buf []T
}

func (s *Stack[T]) Push(el T) {
	s.buf = append(s.buf, el)
}

func (s *Stack[T]) Pop() {
	s.buf = s.buf[:len(s.buf)-1]
}

func (s *Stack[T]) Top() T {
	return s.buf[len(s.buf)-1]
}

func (s *Stack[T]) Len() int {
	return len(s.buf)
}
