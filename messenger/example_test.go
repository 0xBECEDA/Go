package main

import (
	"testing"
	"strconv"
	"io"
	"bufio"
	"os/exec"
)

const(
	MaxClients = 10
)
// RUNSERVER
// запускает сервер и убеждается, что
// он заработал
func runServer(t *testing.T) {
	// готовим на выполнение команду
	// ./server
	cmd := exec.Command("./server")
	// привязываем к ее stsdout пайп
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Errorf("runServer: не удалось привязать пайп к stdout \n")
	}

	r := bufio.NewReader( stdout )
	//запускаем команду
	cmd.Start()
	// проверяем, что напечатал сервак
	string, _ := r.ReadString('\n')

	if string != "SERVER IS ON \n" {
		t.Errorf("Ожидал результат 'SERVER IS ON', а получил %s \n", string)
		t.FailNow()
	}
}

// CHECKMULTIINPUTCLIENT
// проверяет, способен ли клиент отправить самому себе сообщение, состоящиее
// из нескольких строк
func CheckMultiInputClient( stdout io.ReadCloser, stdin io.WriteCloser, id int,
	t *testing.T) {

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
			t.Errorf("checkMultiInputClient id " + idString + ": \n")
			t.Errorf("не удалось вписать строку в буфер \n")
			t.Fail()
		}
	}
	// сливаем в поток
	w.Flush()

	// Убеждаемся, что сообщение пришло от этого же клиента
	string, _ := r.ReadString('\n')
	if string != recieveMsg {
		t.Errorf("checkMultiInputClient id: " + idString)
		t.Errorf("ожидал ' %s ' получил ' %s' ", recieveMsg, string)
		t.Fail()
	}

	// проверяем, что все сообщение дошло в целости и сохранности
	for i := 1; i < 4; i++ {
		string, _ := r.ReadString('\n')

		if string != msg[i] {
			t.Errorf("checkMultiInputClient id: " + idString)
			t.Errorf("ожидал ' %s ' получил ' %s' ", msg[i], string)
			t.Fail()
		}
	}
}

// SendMessageToNonExistendClient
// отправляет сообщение несуществующему клиенту
func SendMessageToNonExistendClient( stdout io.ReadCloser, stdin io.WriteCloser,
	id int, t *testing.T) {
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
			t.Errorf("SendMessageToNonExistendClient id " + MyIdString + ": \n")
			t.Errorf("не удалось вписать строку в буфер \n")
			t.Fail()
		}
	}
	// сливаем в поток
	w.Flush()
	string, _ := r.ReadString('\n')

	if string != recieveMsg {
		t.Errorf("SendMessageToNonExistendClient id: " + MyIdString)
		t.Errorf("ожидал ' %s ' получил ' %s' ", recieveMsg, string)
		t.Fail()
	}
}
// CheckQuitClient
// проверяет, завершается ли клиент, если послать ему
// "quit"
func CheckQuitClient( stdin io.WriteCloser, id int, t *testing.T, cmd *exec.Cmd ) {

	w := bufio.NewWriter( stdin )

	w.WriteString(strconv.Itoa(id) + "\n")
	w.WriteString("quit\n")
	w.Flush()

	// надо как-то проверить, что команда действительно завершается
}

// CheckClientConnectionToServer
// проверяет, подключилсяли клиент к серверу
func CheckClientConnectionToServer( stdout io.ReadCloser,  t *testing.T ) {

	r := bufio.NewReader( stdout )
	string, err := r.ReadString('\n')

	if err == nil {
		if string != "Have a connection with server \n" {
			t.Errorf("CheckClientConnectionToServer connection error: %s \n", string)
			t.FailNow()
		}
	}
}

// RUNCLIENT
// запускает клиента и проводит с ним серию тестов
func ClientTestsGroup(t *testing.T, ch1 chan string, ch2 chan int) {
	//получили id для клиента
	ch1 <- "Get new id"
	id := <- ch2

	cmd := exec.Command("./client", strconv.Itoa(id))
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		t.Errorf("runClient: не удалось привязать пайп к stdout \n")
		t.FailNow()
	}

	stdin, err := cmd.StdinPipe()

	if err != nil {
		t.Errorf("runClient: не удалось привязать пайп к stdin \n")
		t.FailNow()
	}

	cmd.Start()

	t.Run("ClientTestsGroup", func (t *testing.T) {
		t.Run ( "CheckClientConnectionToServer", func (t *testing.T) {
			CheckClientConnectionToServer( stdout, t)
		})

		t.Run ("SendMessageToNonExistendClient", func (t *testing.T) {
			SendMessageToNonExistendClient( stdout, stdin, id, t )
		})

		// t.Logf("id %d: CheckClientConnectionToServer done", id)

		t.Run ("CheckMultiInputClient", func (t *testing.T) {
			CheckMultiInputClient( stdout, stdin, id, t )
		})

		// t.Logf("id %d: CheckMultiInputClient done", id)

		t.Run ("CheckQuitClient", func (t *testing.T) {
			CheckQuitClient( stdin, id, t, cmd )
		})

		// t.Logf("id %d: CheckQuitClient done", id)
	})
}


// createIdForClient
// создает айдишник для клиента, путем увеличения
// глобальной (!) переменной ClientId
func CreateIdForClient(MsgChan chan string, ResultChan chan int) {
	ClientId := 1
	msg := <- MsgChan

	switch msg {

	case "Get new id" :
		ResultChan <- ClientId
		ClientId += 1

	case "Get max id" :
		ResultChan <- ClientId
	}
}

// TESTMESSENGER
// запускет тесты сервера и клиента
func TestMessenger(t *testing.T) {

	CreateIdForClientChanMsg := make( chan string, 100 )
	CreateIdForClientChanResult := make( chan int, 100 )
	go CreateIdForClient( CreateIdForClientChanMsg, CreateIdForClientChanResult )

	t.Run("runServer", func (t *testing.T) {
		runServer( t )
	})

	t.Run("RunClients", func (t *testing.T) {

		// for i := 0; i < 2; i++ {
			t.Run("ClientTestsGroup", func (t *testing.T) {
				t.Parallel()
				ClientTestsGroup( t, CreateIdForClientChanMsg,
					CreateIdForClientChanResult )
			})

			// 	t.Run("ClientTestsGroup2", func (t *testing.T) {
			// 	// t.Parallel()
			// 	ClientTestsGroup( t, CreateIdForClientChanMsg,
			// 		CreateIdForClientChanResult )

			// })
		// }
	})
}
