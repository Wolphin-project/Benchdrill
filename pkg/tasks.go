package exampletasks

import (
	"errors"
	"io/ioutil"
	"net"
	//"fmt"
	//"strconv"
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

	/*addr := net.ParseIP("127.0.0.11")

		if addr == nil {
			return "Invalid IP address", errors.New("tasks: error when parsing string to IP address")
		}

		/*tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.11:6389")

		if err != nil {
			return "Invalid TCP address", err
		}

		service := ":7777"
	    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	    if err != nil {
			return "1", err
		}
	    listener, err := net.ListenTCP("tcp", tcpAddr)
	    if err != nil {
			return "2", err
		}
	    for {
	        conn, err := listener.Accept()
	        if err != nil {
	            continue
	        }
	        go func (conn net.Conn) {
	        	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	    		request := make([]byte, 128) // set maximum request length to 128B to prevent flood based attacks
	    		defer conn.Close()  // close connection before exit
	    		for {
	        		read_len, err := conn.Read(request)

	        		if err != nil {
	            		fmt.Println(err)
	            		break
	        		}

			        if read_len == 0 {
	        		    break // connection already closed by client
	       		 	} else if string(request[:read_len]) == "timestamp" {
	        		    daytime := strconv.FormatInt(time.Now().Unix(), 10)
	       		     	conn.Write([]byte(daytime))
	       		 	} else {
	        		    daytime := time.Now().String()
	       		     	conn.Write([]byte(daytime))
	       		 	}
				}
	    	}(conn)
	   	}


		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7777})

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
		}*/

	return string(result), nil
}
