package write_msg

import (
	"bufio"
	"fmt"
	"messanger/internal"
	"os"
	"strings"
)

func EnterUserName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your user name:")
	name, _ := reader.ReadString('\n')
	return name
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
