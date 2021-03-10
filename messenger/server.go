package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"

	// "io"
)

type sendPackage struct {

	MyID int
	UserID int
	Message []string
	SendStatus int
}

const (
	testUserID = 25
	packSize = 1036
	msgSizeStrings = 10
)

// TODO разберись, как утановить размер буфера, наконец
var buf =  make([]byte, packSize)

func serialization( pack *sendPackage ) ( []byte, error )  {

	buf, err:= json.Marshal( pack )

	if err != nil {
		fmt.Printf(" \n Cound'n serialize data: ", err.Error(), "\n" )
	}
	return buf, err
}

func deserialization( buf []byte ) ( sendPackage, error )  {

	pack := sendPackage{}
	err:= json.Unmarshal( buf, &pack )

	if err != nil {
		fmt.Printf(" \n Cound'n deserialize data: ", err.Error(), "\n" )
	}
	return pack, err
}

func sendMessege( connect *net.TCPConn, len int ) {

	// buf - глобальная var, осторожно!!
	pack, err := deserialization( buf[:len] )

	if err == nil {

		// TODO проверить, зарегестрирован ли юзер,
		// которому отправляется msg
		pack.SendStatus = -1

		sendBuf, err := serialization( &pack )
		len, err := connect.Write( sendBuf )

		if err != nil {
			fmt.Printf("Cann't send: %s \n", err.Error())
		} else {
			fmt.Printf("Bytes sent: %d \n", len)
		}
	}

	return
}

func getMessege( connect *net.TCPConn ) {

	for {
		len, err := connect.Read( buf )

		if err == nil {
			fmt.Printf("message recieved, len %d bytes \n", len );
			sendMessege( connect, len )
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
