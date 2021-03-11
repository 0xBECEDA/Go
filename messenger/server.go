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

type connection struct {

	ClientID int
	Connect *net.TCPConn
}

const (
	testUserID = 25
	packSize = 1000
	msgSizeStrings = 10
	maxClients = 10
)

var ConnectionsTable = make(map[int]*net.TCPConn)

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

func checkErrorSendMessage( err error, len int ) {
	if err != nil {
		fmt.Printf("Cann't send: %s \n", err.Error())
	} else {
		fmt.Printf("Bytes sent: %d \n", len)
	}
	return
}

func sendMessege( myConnect *net.TCPConn, buf []byte, len int ) {

	pack, err := deserialization( buf[:len] )

	if err == nil {
		// проверяем существование юзера, которому
		// отправляем сообщение
		connectUser, found := ConnectionsTable[pack.UserID]
		fmt.Println("found ", found )

		// нашли
		if found == true {
			fmt.Println("connectUser ", connectUser )
			sendBuf, err := serialization( &pack )
			len, err := connectUser.Write( sendBuf )
			checkErrorSendMessage( err, len)

		// не нашли
		} else {
			pack.SendStatus = -1

			sendBuf, err := serialization( &pack )
			len, err := myConnect.Write( sendBuf )
			checkErrorSendMessage( err, len)
		}
	}
	return
}

func getMessege( connect *net.TCPConn, ch chan connection, ch2 chan bool ) {

	remembered := 0
	buf := make([]byte, packSize)

	for {
		len, err := connect.Read( buf )

		if err == nil {
			fmt.Printf("message recieved, len %d bytes \n", len );

			if remembered == 0 {
				pack, err := deserialization( buf[:len] )

				// отправить данные клиента на регистрацию
				if err == nil {
					newConnection:= connection{ ClientID: pack.MyID,
						                        Connect:  connect }
					ch <- newConnection

					// регистрация законцена?
					finished :=  <- ch2
					if finished == true {
						remembered = 1
					}
				}
			}
			sendMessege( connect, buf, len )
		}
	}
	return
}

func RegisterNewClient( ch chan connection, ch2 chan bool ) {

	for {
		newConnection, ok := <- ch

		// канал закрыт?
		if ok == false {
			break
		}
		// сохранить нового клиента
		ConnectionsTable[newConnection.ClientID] = newConnection.Connect

		// проверяем, что действительно сохранилось
		_, found := ConnectionsTable[newConnection.ClientID]

		if found == true {
			fmt.Println("RegisterNewClient: зарегестрирован новый клиент  ",
				newConnection.ClientID,  newConnection.Connect )
			// сигнализируем, что закончили
			ch2 <- found
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

		fmt.Println(" SERVER RUNS \n")
		// создаем поток для регистрации соединений:
		// связываем id клиента с его соединением и запоминаем
		registerClientsChan := make( chan connection, 100 )
		registerClientsChanResult := make( chan bool )
		go RegisterNewClient( registerClientsChan, registerClientsChanResult )

		for {
		conn, err := l.AcceptTCP()

		if err != nil {
			fmt.Println("Accept error", err.Error())
			os.Exit(1)
		}

		// fmt.Printf( "server conn %v \n", conn);
		go getMessege( conn, registerClientsChan, registerClientsChanResult )
		}
	}

	return
}
