package main

import (
	"fmt"
	"net"
	// "io"
	"sync"
	"os"
	"strconv"
	// "unsafe"
	"encoding/json"
	"bufio"
)


type sendPackage struct {

	MyID int
	UserID int
	Message []string
	SendStatus int
}

var myID int

const (
	testUserID = 25
	packSize = 1036
	msgSizeStrings = 10
)

func failedStatus ( userID int ) {

	fmt.Printf("Сообщение пользователю %d не было доставлено\n", userID);
}

func checkSendStatus( buf []byte, len int ) {

	pack, err := deserialization( buf[:len] )

	if err == nil && pack.SendStatus == -1 {
		failedStatus( pack.UserID )
	}
	return
}

//принимает сообщения
func getMessage( connect *net.TCPConn,  wg *sync.WaitGroup ) int {

	// TODO разберись, как уcтановить размер буфера, наконец
	getBuf := make( []byte, packSize )
	for {
		len, err := connect.Read(	getBuf )

		if err == nil {
			// fmt.Printf( "getMessage: error", err.Error(), "\n" )

			checkSendStatus( getBuf, len )
		}
	}
	wg.Done()
	return 0
}


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

func getInput () []string {

	r := bufio.NewReader( os.Stdin )
	input:= make( []string, msgSizeStrings )

	for i := range input {
		string, err := r.ReadString('\n')

		if err != nil {
			fmt.Printf("getInput: Didn't read string \n")
			break

		} else if string == "stop\n" {
			break

		} else {
			input[i] = string
		}
	}

	// for i := range input{
	// 	fmt.Printf("Read: %s \n", input[i])
	// }

	return input
}

func testSerialDeserial() {

	input:= getInput()

	fmt.Scan(&input)
	fmt.Printf("Read: %s \n", input)

	pack := sendPackage{
		MyID: myID,
		UserID: testUserID,
		Message: input,
		SendStatus: 0 }

	fmt.Printf("Msg before serialization: %s \n", pack.Message)

	buf, err := serialization( &pack )

	if err != nil {
		fmt.Printf("Serialization Test failed\n")
	} else {
		pack2, err := deserialization( buf )

		if err != nil {
			fmt.Printf("Deserialization Test failed\n")

		} else {
			fmt.Printf("Msg after serialization %s \n", pack2.Message)
		}
	}
	return
}

// посылает сообщения
func sendMessage( connect *net.TCPConn, wg *sync.WaitGroup ) int {

	pack := sendPackage{
		MyID: myID,
		UserID: testUserID,
		SendStatus: 0 }

	for {
		pack.Message = getInput()
		// fmt.Printf("Read: %s \n", &pack.Message)


		buf, err := serialization( &pack )
		len, err := connect.Write(buf)

		if err != nil {
			fmt.Printf("Cann't send: %s \n", err.Error())
		} else {
			fmt.Printf("Bytes sent: %d \n", len)
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

func Tests () {

	getInput ()
	testSerialDeserial()
}

func main () {

	// Tests ()

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
