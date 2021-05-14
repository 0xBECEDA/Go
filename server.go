package main

import (
    "fmt"
    "net"
    "os"
    "encoding/json"
    "sync"
    "time"
)

type sendPackage struct {

    MyID int
    UserID int
    Message []string
    SendStatus int
}

type connection struct {
    Status string
    ClientID int
    SendQueue chan sendPackage
    Connect *net.TCPConn
}

const (
    packSize = 1000
    msgSizeStrings = 10
)

var ConnectionsTable = make(map[int]connection)

func InitServerState() {
    service := "localhost:3425"
    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    l, err := net.ListenTCP("tcp", tcpAddr)

    if err != nil {
        ExitFailState( err )

    } else {
        RunServerState( l )
    }
}
func ExitFailState( err error ) {
    fmt.Println("--Server error -- :", err.Error())
    os.Exit(1)
}

func RunServerState( l *net.TCPListener ) {
    var wg sync.WaitGroup

    registerClientsChan := make( chan connection, 100 )
    registerClientsChanResult := make( chan bool )

    go RegisterClientsStateMachine( registerClientsChan, registerClientsChanResult, &wg )
    wg.Add(1)

    go AcceptNewConnectionsState( l, registerClientsChan, registerClientsChanResult )
    wg.Wait()
    ExitServerSuccess()
}

func AcceptNewConnectionsState( l *net.TCPListener, ch1 chan connection, ch2 chan bool) {
    fmt.Println("SERVER IS ON \n")
    for {
        conn, err := l.AcceptTCP()

        if err != nil {
            ExitFailState( err )

        } else {
            SupportClientChan := make(chan sendPackage, 100)
            go GetMessage( conn, ch1, ch2, SupportClientChan )
            go SendMessage( conn, ch1,  SupportClientChan )
        }
    }
}
func ExitServerSuccess() {
    os.Exit(0)
}

func RegisterClientsStateMachine( ch chan connection, ch2 chan bool, wg *sync.WaitGroup ) {

    var alive int
    var all int

    for {
        newConnection, ok := <- ch

        if ok == false {
            wg.Done()
            return

        } else {
            state := newConnection.Status

            switch state {

            case "new":
                // сохранить нового клиента
                ConnectionsTable[newConnection.ClientID] = newConnection
                alive++
                all++
                ch2 <- true

            case "dead":
                all--

                if alive > all && alive == 0 {
                    wg.Done()
                    return
                }
            }
        }
    }
    return
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


func GetMessage( connect *net.TCPConn, ch chan connection, ch2 chan bool,
    ch3 chan sendPackage ) {

    RegistrState := "no"
    buf := make([]byte, packSize)

    for {
        len, err := connect.Read( buf )

        // пакет получен
        if err ==  nil {
            fmt.Printf("--GetMessage--: message recieved, len %d bytes \n", len );
            pack, err := deserialization( buf[:len] )

            if err == nil {
                // клиент еще не был зарегистрирован
                if RegistrState == "no" {
                    newConnection:= connection{ Status: "new", ClientID: pack.MyID,
                        SendQueue: ch3, Connect:  connect }
                    ch <- newConnection

                    // регистрация закончена?
                    finished :=  <- ch2
                    if finished == true {
                        RegistrState = "yes"
                    }
                }

                client, found := ConnectionsTable[pack.UserID]
                if found == true {
                    client.SendQueue <- pack

                } else {
                    pack.SendStatus = -1
                    ch3 <- pack
                }
            }
        }
    }
}

func CheckErrorSendMessage( err error, len int ) {
    if err != nil {
        fmt.Printf("--CheckErrorSendMessage--: Can't send: %s \n", err.Error())
    } else {
        fmt.Printf("--CheckErrorSendMessage--: Bytes sent: %d \n", len)
    }
    return
}

func SendMessage( MyConnect *net.TCPConn, ch1 chan connection, ch2 chan sendPackage) {

    for {
        select {
        case pack := <- ch2:

            status := pack.SendStatus
            // проверяем статус отправки
            switch status {
                // клиент сообщил о выходе
            case -5:
                DeadConnect := connection{ Status: "dead", ClientID: pack.MyID }
                ch1 <- DeadConnect
                return

            default:
                sendBuf, err := serialization( &pack )
                len, err := MyConnect.Write( sendBuf )
                CheckErrorSendMessage( err, len)
            }

        default:
            time.Sleep(4 * time.Second)
            pack := sendPackage{ SendStatus: 1 }
            sendBuf, err := serialization( &pack )

            if err == nil {
                len, err := MyConnect.Write( sendBuf )
                CheckErrorSendMessage( err, len)
            }
        }
    }
}


func main () {
    InitServerState()
    return
}
