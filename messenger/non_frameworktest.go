package main

import (
	"strconv"
	"io"
	"bufio"
	"os/exec"
	// "time"
	"fmt"
)

const(
	MaxClients = 10
)

// createIdForClient
// создает айдишник для клиента, путем увеличения
// глобальной (!) переменной ClientId
func CreateIdForClient(MsgChan chan string, ResultChan chan int) {
	ClientId := 1
	for {
		msg := <- MsgChan

		switch msg {

		case "Get new id" :
			ResultChan <- ClientId
			ClientId += 1

		case "Get max id" :
			ResultChan <- ClientId
		}
	}
}

// CompareString
// Читает из заданного стрима строку и сравнивает с ожидаемой строкой.
// Если результатом чтения является ошибка EOF, то пробует считать еще раз -
// скорее всего данные еще дошли в стрим.
// Если и в этом случае ошибка != nil, то возвращает fail
func CompareString(testName string, WaitedString string, r *bufio.Reader, id int ) bool {

	string, err := r.ReadString('\n')

	if err == io.EOF {
		string, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("%s id %d: ошибка чтения: %s \n", testName, id, err )
			return false

		} else if string != WaitedString {
			fmt.Printf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
				id, WaitedString, string )
			return false
		}

	} else if err!= nil {
		fmt.Printf("%s id %d: ошибка чтения: %s \n", testName, id, err )
		return false

	} else if string != WaitedString {
		fmt.Printf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
			id, WaitedString, string )
		return false
	}
	return true
}

// CheckMultiInputClient
// проверяет, способен ли клиент отправить самому себе сообщение, состоящиее
// из нескольких строк
func CheckMultiInputClient( stdout io.ReadCloser, stdin io.WriteCloser, id int ) bool {

	r := bufio.NewReader( stdout )
	w := bufio.NewWriter( stdin )
	idString := strconv.Itoa(id)
	recieveMsg := "Получено сообщение от пользователя " + idString + ":\n"

	// записываем ввод для пакета:
	// - id клиента, которому шлем (себе)
	// - строки, которые шлем
	// - строку "stop\n" - так клиент понимает, что ввод окончен
	// и можно отправлять пакет
	msg := [5]string{ idString + "\n", "Hello me!\n",
		"It is test checkMultiInputClient \n",
		"We will try to print tree strings and recieve them back \n",
		"stop\n" }

	// заливаем все в буфер потока stdin
	for i := 0; i < 5; i++ {
		_, err := w.WriteString( msg[i] )
		if err != nil {
			fmt.Printf("checkMultiInputClient id " + idString + ": \n")
			fmt.Printf("не удалось вписать строку в буфер \n")
			return false
		}
		w.Flush()
	}

	// Убеждаемся, что сообщение пришло от этого же клиента
	retval := CompareString("checkMultiInputClient", recieveMsg, r , id )

	if retval != true {
		return false
	}
	// Проверяем, что все сообщение дошло в целости и сохранности
	for i := 1; i < 4; i++ {
		retval := CompareString("checkMultiInputClient", msg[i], r , id )
		if retval != true {
			return false
		}
	}
	return true
}

// SendMessageToNonExistendClient
// отправляет сообщение несуществующему клиенту
func SendMessageToNonExistendClient( stdout io.ReadCloser, stdin io.WriteCloser,
	id int) bool {
	r := bufio.NewReader( stdout )
	w := bufio.NewWriter( stdin )
	MyIdString := strconv.Itoa(id)
	UserIdString := strconv.Itoa(MaxClients + 1)
	recieveMsg := "Сообщение пользователю c id " + UserIdString + " не было доставлено\n"
	msg := [3]string{ UserIdString + "\n", "Hello other user!\n", "stop\n" }

	// заливаем все в буфер потока stdin
	for i := 0; i < 3; i++ {
		_, err := w.WriteString( msg[i] )
		if err != nil {
			fmt.Printf("SendMessageToNonExistendClient id " + MyIdString + ": \n")
			fmt.Printf("не удалось вписать строку в буфер \n")
			return false
		}
		w.Flush()
	}

	// проверяем, что пришло в ответ от клиента
	return CompareString("SendMessageToNonExistendClient", recieveMsg, r, id )
}

// CheckQuitClient
// проверяет, завершается ли клиент, если послать ему
// "quit"
func CheckQuitClient( stdin io.WriteCloser, id int, cmd *exec.Cmd ) bool {

	w := bufio.NewWriter( stdin )

	w.WriteString(strconv.Itoa(id) + "\n")
	w.WriteString("quit\n")
	w.Flush()

	// надо как-то проверить, что команда действительно завершается
	return true
}

// CheckClientConnectionToServer
// проверяет, подключилсяли клиент к серверу
func CheckClientConnectionToServer( stdout io.ReadCloser, id int ) bool {
	r := bufio.NewReader( stdout )

	msg := "Have a connection with server \n"
	return CompareString("CheckClientConnectionToServer", msg, r , id )
}

// TESTSERVER
// запускает сервер и убеждается, что
// он заработал
func TestServer() bool {
	// готовим на выполнение команду
	// ./server
	cmd := exec.Command("./server")
	// привязываем к ее stsdout пайп
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("TestServer: не удалось привязать пайп к stdout \n")
	}

	r := bufio.NewReader( stdout )
	//запускаем команду
	cmd.Start()
	// проверяем, что напечатал сервак
	string, err := r.ReadString('\n')

	if err != nil {
		fmt.Printf("TestServer: oшибка %s \n", err)
		return false

	} else if string != "SERVER IS ON \n" {
		fmt.Printf("Ожидал результат 'SERVER IS ON', а получил %s \n", string)
		return false
	}
	return true
}

// TestClientGroup
// запускает клиент и проводит с ним серию тестов
func TestClientGroup() bool {

	ch1 := make( chan string, 100 )
	ch2 := make( chan int, 100 )
	go CreateIdForClient( ch1, ch2 )

	for i := 0; i < 2; i++ {

		//получили id для клиента
		ch1 <- "Get new id"
		id := <- ch2

		cmd := exec.Command("./client", strconv.Itoa(id))
		stdout, err1 := cmd.StdoutPipe()
		stdin, err2 := cmd.StdinPipe()

		if err1 != nil {
			fmt.Printf("runClient: не удалось привязать пайп к stdout \n")
			return false
		}

		if err2 != nil {
			fmt.Printf("runClient: не удалось привязать пайп к stdin \n")
			return false
		}

		// запуск клиента
		cmd.Start()

		// проверить наличие соединения с сервером
		retval := CheckClientConnectionToServer( stdout, id)
		if retval == false {
			fmt.Printf("CheckClientConnectionToServer FAILED \n")
			return false

		} else {
			fmt.Printf("CheckClientConnectionToServer PASSED \n")

			// проверить отправку пакета несуществующему клиенту
			retval = SendMessageToNonExistendClient( stdout, stdin, id )
			if retval == false {
				fmt.Printf("SendMessageToNonExistendClient FAILED \n")
				return false

			} else {
				fmt.Printf("SendMessageToNonExistendClient PASSED \n")

				// проверить отправку пакета себе же (в сообщении несколько строк)
				retval = CheckMultiInputClient( stdout, stdin, id )

				if retval == false {
					fmt.Printf("CheckMultiInputClient FAILED \n")
					return false

				} else {
					fmt.Printf("CheckMultiInputClient PASSED \n")

					// завершение клиента
					// retval = CheckQuitClient( stdin, id, cmd )
					// if retval == false {
					// 	fmt.Printf("CheckQuitClient FAILED \n")
					// 	return false
					// }
				}
			}
		}
	}
	return true
}


func main () {

	retval := TestServer()
	retval2 := TestClientGroup()
	if retval == false || retval2 == false {
		// if retval2 == false {
		fmt.Printf("Test Messenger FAILED \n")
	}
}
