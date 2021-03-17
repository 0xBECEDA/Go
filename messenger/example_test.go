package main

import (
	"testing"
	"strconv"
	"io"
	"bufio"
	"os/exec"
	// "time"
	// "fmt"
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
func CompareString(t *testing.T, testName string, WaitedString string,
	stdout io.ReadCloser, id int ) {

	r := bufio.NewReader( stdout )
	string, err := r.ReadString('\n')

	if err == io.EOF {
		string, err := r.ReadString('\n')
		if err != nil {
			t.Errorf("%s id %d: ошибка чтения: %s \n", testName, id, err )
			t.FailNow()

		} else if string != WaitedString {
			t.Errorf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
				id, WaitedString, string )
			t.FailNow()
		}

	} else if err!= nil {
		t.Errorf("%s id %d: ошибка чтения: %s \n", testName, id, err )
		t.FailNow()

	} else if string != WaitedString {
		t.Errorf("%s id %d: ожидал строку: \n %s а получил: \n %s \n", testName,
			id, WaitedString, string )
		t.FailNow()
	}
}

// CHECKMULTIINPUTCLIENT
// проверяет, способен ли клиент отправить самому себе сообщение, состоящиее
// из нескольких строк
func CheckMultiInputClient( stdout io.ReadCloser, stdin io.WriteCloser, id int,
	t *testing.T) {

	// r := bufio.NewReader( stdout )
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
	CompareString(t , "checkMultiInputClient", recieveMsg, stdout , id )

	// Проверяем, что все сообщение дошло в целости и сохранности
	for i := 1; i < 4; i++ {
		CompareString(t , "checkMultiInputClient", msg[i], stdout , id )
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
			t.FailNow()
		}
	}
	// сливаем в поток
	w.Flush()

	string, err := r.ReadString('\n')

	// проверям, что пришло в ответ от клиента
	CompareString(t , "SendMessageToNonExistendClient", recieveMsg, stdout , id )
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
func CheckClientConnectionToServer( stdout io.ReadCloser,  t *testing.T, id int ) {

	msg := "Have a connection with server \n"
	CompareString(t , "CheckClientConnectionToServer", msg, stdout , id )
}

// TESTSERVER
// запускает сервер и убеждается, что
// он заработал
func TestServer(t *testing.T) {
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
	string, err := r.ReadString('\n')

	if err != nil {
		t.Errorf("runServer: oшибка %s \n", err)

	} else if string != "SERVER IS ON \n" {
		t.Errorf("Ожидал результат 'SERVER IS ON', а получил %s \n", string)
		t.FailNow()
	}
}

// TestClientGroup
// запускает клиент и проводит с ним серию тестов
func TestClientGroup(t *testing.T) {

	ch1 := make( chan string, 100 )
	ch2 := make( chan int, 100 )
	go CreateIdForClient( ch1, ch2 )

	// for i := 0; i < 2; i++ {

	//получили id для клиента
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

	cmd.Start()
	t.Logf("\n id %d: Run all tests", id)

	// запускает тесты над открытым клиентом
	t.Run("ClientTestsGroup", func (t *testing.T) {
		t.Run ( "CheckClientConnectionToServer", func (t *testing.T) {
			CheckClientConnectionToServer( stdout, t, id)
		})

		t.Run ("SendMessageToNonExistendClient", func (t *testing.T) {
			SendMessageToNonExistendClient( stdout, stdin, id, t )
		})

		t.Run ("CheckMultiInputClient", func (t *testing.T) {
			CheckMultiInputClient( stdout, stdin, id, t )
		})

		t.Run ("CheckQuitClient", func (t *testing.T) {
			CheckQuitClient( stdin, id, t, cmd )
		})

	})
	// }
}
