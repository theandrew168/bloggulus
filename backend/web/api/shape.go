package api

import (
	"fmt"
	"net/url"
	"strconv"
)

const MinPage = 1
const DefaultPage = 1

const MinSize = 1
const MaxSize = 50
const DefaultSize = 20

type ItemsResponse[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type ItemResponse[T any] struct {
	Item T `json:"item"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var ErrInvalidPage = fmt.Errorf("invalid page (should be an integer greater than or equal to %d)", MinPage)
var ErrPageOutOfRange = fmt.Errorf("page out of range (should be greater than or equal to %d)", MinPage)

func ParsePage(values url.Values) (int, error) {
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		return 0, ErrInvalidPage
	}

	if page < MinPage {
		return 0, ErrPageOutOfRange
	}

	return page, nil
}

var ErrInvalidSize = fmt.Errorf("invalid size (should be an integer between %d and %d)", MinSize, MaxSize)
var ErrSizeOutOfRange = fmt.Errorf("size out of range (should be between %d and %d)", MinSize, MaxSize)

func ParseSize(values url.Values) (int, error) {
	size, err := strconv.Atoi(values.Get("size"))
	if err != nil {
		return 0, ErrInvalidSize
	}

	if size < MinSize {
		return 0, ErrSizeOutOfRange
	}
	if size > MaxSize {
		return 0, ErrSizeOutOfRange
	}

	return size, nil
}
