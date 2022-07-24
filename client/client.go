package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type sendPackage struct {
	MyID       int
	UserID     int
	Message    []string
	SendStatus int
}

var myID int

const (
	packSize       = 1036
	msgSizeStrings = 10
)

var ChatTable = make(map[int]bool)
var alive int
var new int

func GetClientId() error {

	id, err := strconv.Atoi(os.Args[1])

	if err != nil {
		return err

	} else {
		myID = id
	}
	return nil
}

func ConnectToServer() (*net.TCPConn, error) {
	serVaddr := "localhost:8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serVaddr)

	if err != nil {
		return nil, err

	} else {
		connect, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			fmt.Printf("Connection failed: ", err.Error(), "\n")
			return nil, err
		}
		return connect, nil
	}
}

func checkSendStatus(SendStatus int, UserID int) int {

	switch SendStatus {

	case -1:
		failedStatus(UserID)
		return -1
	case 1:
		return 1
	default:
		return 0
	}
}
func failedStatus(userID int) {

	fmt.Printf("Сообщение пользователю c id %d не было доставлено\n", userID)
}
func printRecievedMessage(UserID int, buf []string) {
	fmt.Printf("Получено сообщение от пользователя %d:\n", UserID)

	for i := range buf {

		if buf[i] == "" {
			break
		}
		fmt.Printf(buf[i])
	}
	return
}

//принимает сообщения
// func getMessage( connect *net.TCPConn,  wg *sync.WaitGroup ) int {
func getMessage(connect *net.TCPConn, ch chan string) int {
	getBuf := make([]byte, packSize)

	for {
		timeoutDuration := 10 * time.Second
		connect.SetReadDeadline(time.Now().Add(timeoutDuration))
		len, err := connect.Read(getBuf)

		if err == nil {
			pack, err := deserialization(getBuf[:len])

			if err == nil && 0 == checkSendStatus(pack.SendStatus, pack.UserID) {
				printRecievedMessage(pack.UserID, pack.Message)
			}

		} else {
			fmt.Printf("Ошибка чтения: возможно разорвано соединение\n")
			break
		}
	}
	ch <- "quit"
	// wg.Done()
	return 0
}

func countDialogs(ch chan string) {
	for {
		msg := <-ch
		if msg == "new" {
			alive += 1
			new += 1
			// fmt.Printf("countDialogs: amount of chats %d \n", new)
		} else if msg == "quit" {
			alive -= 1
			// fmt.Printf("countDialogs: amount of alive chats %d \n", alive)
		}
	}
}

func IsAnyChatAlive(wg *sync.WaitGroup) {
	for {
		if new > 0 && alive == 0 {
			// fmt.Printf("IsAnyChatAlive done \n ")
			wg.Done()
			break
		}
	}
}

// посылает сообщения
// func sendMessage( connect *net.TCPConn, wg *sync.WaitGroup, ch chan SendPackage ) int {

func sendMessage(connect *net.TCPConn, ch chan sendPackage, ch2 chan bool) int {
	for {
		pack := <-ch
		buf, err := serialization(&pack)

		if err == nil {
			_, err := connect.Write(buf)

			if err != nil {
				fmt.Printf("Cann't send: %s \n", err.Error())
			} else {
				// fmt.Printf("Bytes sent: %d \n", len)
				if pack.SendStatus == -5 {
					ch2 <- true
				}
			}
		}
	}
	// wg.Done()
	return 0
}
func serialization(pack *sendPackage) ([]byte, error) {

	buf, err := json.Marshal(pack)

	if err != nil {
		fmt.Printf(" \n Cound'n serialize data: ", err.Error(), "\n")
	}
	return buf, err
}
func deserialization(buf []byte) (sendPackage, error) {

	pack := sendPackage{}
	err := json.Unmarshal(buf, &pack)

	if err != nil {
		fmt.Printf(" \n Cound'n deserialize data: ", err.Error(), "\n")
	}
	return pack, err
}

func IsItNewChat(id int) bool {
	_, found := ChatTable[id]
	return found
}
func getInput() ([]string, int) {

	r := bufio.NewReader(os.Stdin)
	input := make([]string, msgSizeStrings)
	inputStatus := 0

	for i := range input {
		string, err := r.ReadString('\n')

		if err != nil {
			fmt.Printf("getInput: Didn't read string \n")
			break

		} else if string == "stop\n" {
			break

		} else if string == "quit\n" {
			inputStatus = -5
			break
		} else {
			input[i] = string
		}
	}

	// for i := range input{
	// 	fmt.Printf("Read: %s \n", input[i])
	// }

	return input, inputStatus
}

func makePackage(input []string, clientStatus int) (sendPackage, error) {

	// убираем '\n' после числа
	string := input[0]
	n := len(string) - 1

	RecieverUserID, err := strconv.Atoi(string[:n])
	pack := sendPackage{}
	if err != nil {
		fmt.Printf("makePackage: Didn't convert id of user: ", err.Error(), "\n")

	} else {
		pack.MyID = myID
		pack.UserID = RecieverUserID
		pack.Message = input[1:]
		pack.SendStatus = clientStatus

	}
	return pack, err
}
func testSerialDeserial() {

	input, status := getInput()

	fmt.Scan(&input)
	fmt.Printf("Read: %s \n", input)

	pack, err := makePackage(input, status)
	if err == nil {
		fmt.Printf("Msg before serialization: %s \n", pack.Message)

		buf, err := serialization(&pack)

		if err != nil {
			fmt.Printf("Serialization Test failed\n")
		} else {
			pack2, err := deserialization(buf)

			if err != nil {
				fmt.Printf("Deserialization Test failed\n")

			} else {
				fmt.Printf("Msg after serialization %s \n", pack2.Message)
			}
		}
	}
	return
}

func Tests() {

	getInput()
	testSerialDeserial()
}
func driverLoop(wg *sync.WaitGroup, DoneChannel chan string, connect *net.TCPConn) {
	SendPackageQueue := make(chan sendPackage, 100)
	SendLastPackageResult := make(chan bool, 100)

	for {
		input, clientStatus := getInput()
		SendPack, err := makePackage(input, clientStatus)

		if err == nil {
			found := IsItNewChat(SendPack.UserID)

			if found == false {
				go sendMessage(connect, SendPackageQueue, SendLastPackageResult)
				go getMessage(connect, DoneChannel)
				DoneChannel <- "new"
				// go sendMessage( connect, &wg, SendPackageQueue )
				// wg.Add(1)
				// go getMessage( connect, &wg )
				// wg.Add(1)
				// wg.Wait()
			}
			SendPackageQueue <- SendPack
		}
		// клиент решил выйти из чата
		if clientStatus == -5 {
			// удостовериваемся, что последний отправленный пакет дошел
			result := <-SendLastPackageResult
			if result == true {
				// выключаем поток IsAnyChatAlive (его отслеживает main,
				// чтоб выйти )
				wg.Done()
				break
			}
		}
	}
}

/*func main() {

	// Tests ()

	err := GetClientId()

	if err != nil {
		fmt.Printf(" Didn't get client id: ", err.Error(), "\n")
		os.Exit(1)
	}

	connect, err := ConnectToServer()

	if err != nil {
		fmt.Printf("\n Connection failed, exit \n")
		os.Exit(1)
	}

	fmt.Printf("Have a connection with http \n")
	var wg sync.WaitGroup
	// отслеживаем, сколько "живых" диалогов
	DoneChannel := make(chan string, 100)
	go countDialogs(DoneChannel)
	go IsAnyChatAlive(&wg)
	wg.Add(1)
	go driverLoop(&wg, DoneChannel, connect)
	wg.Wait()
}
*/
