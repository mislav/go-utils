package utils

type Set struct {
	data map[string]bool
}

func NewSet() *Set {
	return &Set{make(map[string]bool)}
}

func NewSetFromSlice(values []string) *Set {
	set := NewSet()
	for _, value := range values {
		set.Add(value)
	}
	return set
}

func (s *Set) Add(value string) {
	s.data[value] = true
}

func (s *Set) Remove(value string) {
	delete(s.data, value)
}

func (s *Set) Values() []string {
	values := make([]string, len(s.data))
	i := 0
	for value, _ := range s.data {
		values[i] = value
		i += 1
	}
	return values
}

func (s *Set) Length() int {
	return len(s.data)
}
