package main

import (
	"fmt"
	"net"
	"io"
	"sync"
	"os"
	"strconv"
	// "mem"
	"unsafe"
	// "bufio"
)


type sendPackage struct {

	myID int
	userID int
    message [msgSize]byte
	send_status int

}

var myID int

const (
	testUserID = 25
	packSize = 1036
	msgSize = 1024

)

func failedStatus ( userID int ) {

	fmt.Printf("Сообщение пользователю %d не было доставлено\n", userID);

}

//принимает сообщения
func getMessage( connect *net.TCPConn,  wg *sync.WaitGroup ) int {

	buf := make( []byte, packSize )
	for {
		_, err := connect.Read(buf)

		if err == io.EOF {
			break
		}
	}
	wg.Done()
	return 0
}

// спизжено из https://github.com/freboat/gomem/blob/master/mem/mem.go

type usp unsafe.Pointer
type size_t int

func Memcpy(dest, src unsafe.Pointer, len size_t) unsafe.Pointer {

	cnt := len >> 3
	var i size_t = 0
	for i = 0; i < cnt; i++ {
		var pdest *uint64 = (*uint64)(usp(uintptr(dest) + uintptr(8*i)))
		var psrc *uint64 = (*uint64)(usp(uintptr(src) + uintptr(8*i)))
		*pdest = *psrc
	}
	left := len & 7
	for i = 0; i < left; i++ {
		var pdest *uint8 = (*uint8)(usp(uintptr(dest) + uintptr(8*cnt+i)))
		var psrc *uint8 = (*uint8)(usp(uintptr(src) + uintptr(8*cnt+i)))

		*pdest = *psrc
	}
	return dest
}


func serialization( pack *sendPackage, buf *[]byte )  {

	return
}

// посылает сообщения
func sendMessage( connect *net.TCPConn, wg *sync.WaitGroup ) int {

	buf := make( []byte, packSize )
	var input [msgSize]byte
	pack := sendPackage{myID: myID, userID: testUserID, message: input, send_status: 0}
	for {
		fmt.Scan(&pack.message)
		fmt.Printf("Read: %s \n", &pack.message);

		//todo serialization
		len, err := connect.Write(buf)

		if err != nil {
			fmt.Printf("Cann't send: %s \n", err.Error());
		} else {
			fmt.Printf("Bytes sent: %d \n", len);
		}
	}

	wg.Done()
	return 0
}

func ConnectToServer() ( *net.TCPConn, error ) {

	serVaddr :=  "localhost:3425"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serVaddr)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		return nil, err

	} else {

		// fmt.Printf( "Client tcp addr: %v \n",  tcpAddr)
		connect , err := net.DialTCP( "tcp", nil, tcpAddr)

		if err != nil {
			fmt.Printf( "Connection failed: ", err.Error(), "\n" )
			return nil, err

		}
		return connect, nil
	}

}

func GetClientId() error {

	id, err := strconv.Atoi( os.Args[1] )

	if err != nil {
		return err

	} else {
		myID = id
	}
	return nil
}

func main () {

	err := GetClientId()

	if err != nil {
		fmt.Printf( " Didn't get client id: ", err.Error(), "\n" )
		os. Exit(1)
	}

	fmt.Printf( " Client id %d \n", myID )

	connect, err := ConnectToServer()

	if err != nil {
		fmt.Printf( "\n Connection failed, exit \n" )
		os. Exit(1)
	}

	fmt.Printf( " Have a connection with server \n" )
	var wg sync.WaitGroup

	go sendMessage( connect, &wg )
	wg.Add(1)
	go getMessage( connect, &wg )
	wg.Add(1)
	wg.Wait()

}
