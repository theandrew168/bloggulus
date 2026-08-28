package value

import "errors"

var ErrNegativeCount = errors.New("count: value cannot be negative")

type Count struct {
	value int
}

func NewCount(value int) (Count, error) {
	if value < 0 {
		return Count{}, ErrNegativeCount
	}

	c := Count{
		value: value,
	}
	return c, nil
}

func (c Count) Value() int {
	return c.value
}

func (c Count) Increment() Count {
	return Count{
		value: c.Value() + 1,
	}
}
