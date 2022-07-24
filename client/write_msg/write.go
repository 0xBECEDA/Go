package write_msg

import (
	"bufio"
	"fmt"
	"messanger/internal"
	"os"
	"strings"
)

func EnterHost() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your host for yor server:")
	host, _ := reader.ReadString('\n')
	return strings.Trim(host, "\n")
}

func Authorize(host string) *internal.AuthorizeMessage {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your user name:")
	name, _ := reader.ReadString('\n')

	return &internal.AuthorizeMessage{
		Name: strings.Trim(name, "\n"),
		Host: host,
	}
}

func GetInput(ch chan internal.Message, myUserName string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter the name of the user you want to have a dialogue with:")
		name, _ := reader.ReadString('\n')

		fmt.Println("Enter message:")
		text, _ := reader.ReadString('\n')

		ch <- internal.Message{
			FromName: myUserName,
			ToName:   strings.Trim(name, "\n"),
			Data:     []byte(text),
		}
	}
}
