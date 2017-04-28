package exampletasks

import (
	"errors"
	"io/ioutil"
	"net"
	"os/exec"
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
	service := "one-mega-nginx:80"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

	if err != nil {
		return "Invalid TCP address", err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return "Error when establishing TCP connection", err
	}

	_, err = conn.Write([]byte("GET /one_mega_file HTTP/1.0\r\n\r\n"))

	if err != nil {
		return "Error when writing into the connection", err
	}

	result, err := ioutil.ReadAll(conn)

	if err != nil {
		return "Error when reading from the TCP socket", err
	}

	return string(result), nil
}

// Sysbench
func SysbenchTask(args ...string) (string, error) {
	cmd := "sysbench"

	if err := exec.Command(cmd, args...).Run(); err != nil {
		return "Error when executing sysbench", err
	}

	return "Sysbench launched successfully", nil
}
