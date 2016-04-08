package main

import (
	"errors"
	"strings"
)

// StringService represents an object that will implement the StringService
// interface
type StringService struct{}

// Uppercase implements StringService
func (StringService) Uppercase(str string) (string, error) {
	return strings.ToUpper(str), nil
}

// Count implements StringService
func (StringService) Count(str string) int {
	return len(str)
}
