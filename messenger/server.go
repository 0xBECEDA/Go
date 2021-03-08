package main

import (
	"fmt"
	"net"
	"os"
	// "io"
)

func getMessege( connect *net.TCPConn ) {

	buf := make([]byte, 1024)

	for {
		_, err := connect.Read(buf)

		if err == nil {
			fmt.Printf("message recieved \n");
		}
	}
	return
}

func main () {

	service := "localhost:3425"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	l, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Println("Server: listening error", err.Error())
		os.Exit(1)

	} else {
		fmt.Println("Server: l %v", l)
		for {
		conn, err := l.AcceptTCP()

		if err != nil {
			fmt.Println("Accept error", err.Error())
			os.Exit(1)
		}

		fmt.Printf( "server conn %v \n", conn);
		go getMessege( conn )
		}
	}

	return
}
