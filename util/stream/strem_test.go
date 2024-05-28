package stream

import "testing"

func TestStream_Max(t *testing.T) {
	type A struct {
		Name int
	}

	s := New(A{1}, A{2}, A{3})
	t.Log(s.Max("Name"))
}
