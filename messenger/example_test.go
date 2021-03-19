package main

import (
	"testing"
	"strconv"
	"io"
	"bufio"
	"os/exec"
	"sync"
	// "time"
	// "fmt"
)

const(
	MaxClients = 2
)

type clientData struct {
	ID int
	Stdin io.WriteCloser
	Stdout io.ReadCloser
}

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
func CompareString(t *testing.T, testName string, WaitedStrings[]string, r *bufio.Reader,
	id int ) {

	for i := range WaitedStrings {

		string, err := r.ReadString('\n')

		if err == io.EOF {
			string, err := r.ReadString('\n')
			if err != nil {
				t.Errorf("%s id %d: ошибка чтения: %s \n", testName, id, err )
				t.Fail()

			} else if string != WaitedStrings[i] {
				t.Errorf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
					id, WaitedStrings[i], string )
				t.Fail()
			}

		} else if err!= nil {
			t.Errorf("%s id %d: ошибка чтения: %s \n", testName, id, err )
			t.Fail()

		} else if string != WaitedStrings[i] {
			t.Errorf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
				id,  WaitedStrings[i], string )
			t.Fail()
		}
	}
}

// SendMultiStringMsg
// Отправляет сообщение из нескольких строк клиенту, по заданному ID и проверяет
// вывод
func SendMultiStringMsg( stdout io.ReadCloser, stdin io.WriteCloser, id int,
	t *testing.T) {

	r := bufio.NewReader( stdout )
	w := bufio.NewWriter( stdin )
	idString := strconv.Itoa(id)
	recieveMsg := make([]string, 1 )
	recieveMsg[0] = "Получено сообщение от пользователя " + idString + ":\n"

	// записываем ввод для пакета:
	// - id клиента, которому шлем
	// - строки, которые шлем
	// - строку "stop\n" - так клиент понимает, что ввод окончен
	// и можно отправлять пакет
	msg := [5]string{ idString + "\n", "Hello!\n",
		"It is test SendMultiStringMsg from client " + idString + "\n",
		"We will try to print tree strings and recieve them back \n",
		"stop\n" }

	// заливаем все в буфер потока stdin
	for i := 0; i < 5; i++ {
		_, err := w.WriteString( msg[i] )
		if err != nil {
			t.Errorf("SendMultiStringMsg id " + idString + ": \n")
			t.Errorf("не удалось вписать строку в буфер \n")
			t.Fail()
		}
	}
	// сливаем в поток
	w.Flush()

	// Убеждаемся, что сообщение пришло от того клиента, которому посылали
	CompareString(t , "SendMultiStringMsg", recieveMsg, r, id )

	// Проверяем, что все сообщение дошло в целости и сохранности
	CompareString(t , "SendMultiStringMsg", msg[ 1 : 4 ], r, id )
}

// SendMessageToNonExistendClient
// отправляет сообщение несуществующему клиенту
func SendMessageToNonExistendClient( stdout io.ReadCloser, stdin io.WriteCloser,
	id int, t *testing.T) {
	r := bufio.NewReader( stdout )
	w := bufio.NewWriter( stdin )
	MyIdString := strconv.Itoa(id)
	UserIdString := strconv.Itoa(MaxClients + 1)
	recieveMsg := make([]string, 1)
	recieveMsg[0] = "Сообщение пользователю c id " + UserIdString + " не было доставлено\n"
	msg := [3]string{ UserIdString + "\n", "Hello other user!\n", "stop\n" }

	// заливаем все в буфер потока stdin
	for i := 0; i < 3; i++ {
		_, err := w.WriteString( msg[i] )
		if err != nil {
			t.Errorf("SendMessageToNonExistendClient id " + MyIdString + ": \n")
			t.Errorf("не удалось вписать строку в буфер \n")
			t.Fail()
		}
	}
	// сливаем в поток
	w.Flush()
	// проверям, что пришло в ответ от клиента
	CompareString(t , "SendMessageToNonExistendClient", recieveMsg, r, id )
}

// CheckQuitClient
// проверяет, завершается ли клиент, если послать ему
// "quit"
func CheckQuitClient( stdin io.WriteCloser, id int, t *testing.T, wg *sync.WaitGroup ) {

	w := bufio.NewWriter( stdin )

	w.WriteString(strconv.Itoa(id) + "\n")
	w.WriteString("quit\n")
	w.Flush()
	wg.Done()
	// надо как-то проверить, что команда действительно завершается
}

// CheckClientConnectionToServer
// проверяет, подключилсяли клиент к серверу
func CheckClientConnectionToServer( stdout io.ReadCloser,  t *testing.T, id int,
	wg *sync.WaitGroup) {
	r := bufio.NewReader( stdout )
	msg := make([]string, 1)
	msg[0] = "Have a connection with server \n"
	CompareString(t , "CheckClientConnectionToServer", msg, r, id )
	wg.Done()
}

// TESTSERVER
// запускает сервер и убеждается, что
// он заработал
func RunServer(t *testing.T) {
	// готовим на выполнение команду
	// ./server
	cmd := exec.Command("./server")
	// привязываем к ее stsdout пайп
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Errorf("runServer: не удалось привязать пайп к stdout \n")
		t.FailNow()
	}

	r := bufio.NewReader( stdout )
	//запускаем команду
	cmd.Start()
	// проверяем, что напечатал сервак
	string, err := r.ReadString('\n')

	if err != nil {
		t.Errorf("runServer: oшибка %s \n", err)

	} else if string != "SERVER IS ON \n" {
		t.Errorf("Ожидал результат 'SERVER IS ON', а получил %s \n", string)
		t.FailNow()
	}
}

// ClientCommunication
// Отсылает всем клиентам в массиве сообщения от клиента, который находится
// по заданному индексу в массе
func ClientCommunication( t *testing.T, wg *sync.WaitGroup, indx int, clients []clientData) {

	curClient := clients[indx]

	for i := range clients {
		otherClient := clients[i]
		SendMultiStringMsg( otherClient.Stdout, curClient.Stdin, otherClient.ID, t )
	}
	wg.Done()
}

// ClientGroupTest
// Запускает заданное количество клиентов и инициирует "общение", т.е.
// передачу сообщений между ними: каждый клиент должен отправить
// каждому сообщеие и получить
func ClientGroupTest( t *testing.T ) {

	ch1 := make( chan string, 100 )
	ch2 := make( chan int, 100 )
	go CreateIdForClient( ch1, ch2 )

	clients := make([]clientData, MaxClients)
	var wg sync.WaitGroup

	for i := 0; i < MaxClients; i++ {

		ch1 <- "Get new id"
		id := <- ch2

		cmd := exec.Command("./client", strconv.Itoa(id))
		stdout, err1 := cmd.StdoutPipe()
		stdin, err2 := cmd.StdinPipe()

		if err1 != nil {
			t.Errorf("runClient: не удалось привязать пайп к stdout \n")
			t.FailNow()
		}

		if err2 != nil {
			t.Errorf("runClient: не удалось привязать пайп к stdin \n")
			t.FailNow()
		}

		client := clientData{ ID: id, Stdin: stdin, Stdout: stdout }
		clients[i] = client
		cmd.Start()
	}

	// проверяем, что все клиенты подключены к серверу
	t.Run("CheckClientConnectionToServer", func ( t *testing.T ) {
		for i := range clients {
			go CheckClientConnectionToServer( clients[i].Stdout, t, clients[i].ID, &wg)
			wg.Add(1)
		}
		wg.Wait()
	})

	// каждый клиент отправляет самому себе сообщение - так сервер его зарегистрирует
	t.Run("SendMultiStringMsg", func ( t *testing.T ) {
		for i := range clients {
			SendMultiStringMsg( clients[i].Stdout, clients[i].Stdin, clients[i].ID, t)
		}
	})

	// каждый клиент отправляет сообщение несуществующему клиенту
	t.Run("SendMessageToNonExistendClient", func ( t *testing.T ) {
		for i := range clients {
			SendMessageToNonExistendClient( clients[i].Stdout, clients[i].Stdin,
				clients[i].ID, t)
		}
	})

	// каждый клиент отправляет другому клиенту сообщение
	t.Run("ClientCommunication", func ( t *testing.T ) {
		for i := range clients {
			go ClientCommunication( t, &wg, i, clients)
			wg.Add(1)
		}
		wg.Wait()
	})

	// каждый клиент выходит и сообщает это серверу
	t.Run("CheckQuitClient", func ( t *testing.T ) {
		for i := range clients {
			go CheckQuitClient( clients[i].Stdin, clients[i].ID, t, &wg)
			wg.Add(1)
		}
	wg.Wait()
	})
}



func TestMessenger(t *testing.T) {

	t.Run("TestServer", func (t *testing.T) {
		RunServer(t)
	})

	t.Run("ClientGroupTest", func (t *testing.T) {
		ClientGroupTest(t)
	})
}
