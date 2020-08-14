package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/grpc-demo/proto"
)

func RunTcp() {
	// 客户端
	conn, err := net.Dial("tcp", "127.0.0.1:3200")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n')
		message := strings.Trim(input, "\r\n")
		if strings.ToUpper(message) == "Q" {
			return
		}
		data, err := proto.Encode(message)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("err :", err)
			return
		}
	}
}
