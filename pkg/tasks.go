package exampletasks

import (
	"time"
	"runtime"
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

// Sleeping
func RestfulSleep() (string, error) {
	time.Sleep(1000 * time.Millisecond)

	return "Slept 1 s.", nil
}

// Busy
func GetBusy() (string, error) {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
					case <-done:
						return
					default:
				}
			}
		}()
	}

	time.Sleep(time.Second * 5)
	close(done)

	return "Hard work done.", nil
}
