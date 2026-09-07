package value

type Set[T comparable] struct {
	values map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] {
	s := Set[T]{
		values: make(map[T]struct{}),
	}
	for _, value := range values {
		s.Add(value)
	}
	return &s
}

func (s *Set[T]) Add(value T) {
	s.values[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.values, value)
}

func (s *Set[T]) Contains(value T) bool {
	_, ok := s.values[value]
	return ok
}

func (s *Set[T]) Values() []T {
	values := make([]T, 0, len(s.values))
	for value := range s.values {
		values = append(values, value)
	}
	return values
}

func (s *Set[T]) Len() int {
	return len(s.values)
}

// Returns a new set of values that are in both the current set and the other set.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	intersect := NewSet[T]()
	for value := range s.values {
		if other.Contains(value) {
			intersect.Add(value)
		}
	}
	return intersect
}

// Returns a new set of values that are in the current set or the other set.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	union := NewSet[T]()
	for value := range s.values {
		union.Add(value)
	}
	for value := range other.values {
		union.Add(value)
	}
	return union
}

// Returns a new set of values that are in the current set but not in the other set.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	difference := NewSet[T]()
	for value := range s.values {
		if !other.Contains(value) {
			difference.Add(value)
		}
	}
	return difference
}
