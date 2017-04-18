package exampletasks

import (
	"errors"
	"io/ioutil"
	"net"
	"runtime"
	"time"
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

// TCP Socket
func OperateTCP() (string, error) {
	addr := net.ParseIP("127.0.0.11")

	if addr == nil {
		return "Invalid IP address", errors.New("tasks: error when parsing string to IP address")
	}

	/*tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.11:6389")

	if err != nil {
		return "Invalid TCP address", err
	}*/

	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 6389})

	if err != nil {
		return "Error when establishing TCP connection", err
	}

	_, err = conn.Write([]byte("You see that awesome dial?"))

	if err != nil {
		return "Error when writing into the connection", err
	}

	result, err := ioutil.ReadAll(conn)

	if err != nil {
		return "Error when reading from the TCP socket", err
	}

	return string(result), nil
}
