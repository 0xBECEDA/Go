package main

import (
	"fmt"
	"net"
	"io"
	"sync"
	"os"
	// "bufio"
)

var buf []byte

//принимает сообщения
func getMessage( connect *net.TCPConn,  wg *sync.WaitGroup ) int {

	for {
		_, err := connect.Read(buf)

		if err == io.EOF {
			break
		}
	}
	wg.Done()
	return 0
}

// func getInput() []byte {

// }

// посылает сообщения
func sendMessage( connect *net.TCPConn, wg *sync.WaitGroup ) int {

	var input []byte

	for {
		fmt.Scan(&input)
		fmt.Printf("Read: %s \n", input);
		len, err := connect.Write(input)

		if err != nil {
			fmt.Printf("Cann't send: %s \n", err.Error());
		} else {
			fmt.Printf("Bytes sent: %d \n", len);
		}
	}

	wg.Done()
	return 0
}


func main () {

	// установить соединение
	serVaddr :=  "localhost:3425"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serVaddr)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	fmt.Printf( "Client tcp addr: %v \n",  tcpAddr);
	connect , err := net.DialTCP( "tcp", nil, tcpAddr)

	if err != nil {
		fmt.Printf( "Connection failed: ", err.Error(), "\n" );

	} else {
		fmt.Printf( "Have connection \n" );
		fmt.Printf( "%v \n", connect);

		var wg sync.WaitGroup

		go sendMessage( connect, &wg )
		wg.Add(1)
		go getMessage( connect, &wg )
		wg.Add(1)
		wg.Wait()
	}
}
