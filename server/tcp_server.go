package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/grpc-demo/proto"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		receverStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", receverStr)
		conn.Write([]byte(receverStr)) // 发送数据
	}

}

func decodeProcess(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)
		data, err := proto.Encode(msg)
		if err != nil {
			fmt.Println("receive msg  resend, err:", err)
			return
		}
		conn.Write(data) // 发送数据
	}

}

func RunTcp() {
	listen, err := net.Listen("tcp", "127.0.0.1:3200")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		// go process(conn)
		go decodeProcess(conn)
	}
}
