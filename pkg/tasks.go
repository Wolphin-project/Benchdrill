package exampletasks

import (
	"errors"
)

// Add ...
func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

// Multiply ...
func Multiply(args ...int64) (int64, error) {
	sum := int64(1)
	for _, arg := range args {
		sum *= arg
	}
	return sum, nil
}

// PanicTask ...
func PanicTask() (string, error) {
	panic(errors.New("oops"))
}

// Un petit test
func SimpleTest() (string, error) {
	return "Quel test exceptionnel", nil
}
