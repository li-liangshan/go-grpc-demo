package client

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-demo/message"
	"google.golang.org/grpc"
)

type Hello struct {
	Name   string
	Age    int32
	School string
}

type What struct {
	Id int64
	Hello
}

func Start() {
	what := What{Id: 1}
	what.Hello = Hello{Name: "xiaoming", Age: 23, School: "what school"}
	fmt.Println("hello id: ", what.Id)
	fmt.Println("what hello name: ", what.Hello.Name)
	what.Hello.Name = "xiaoming2"
	what.Hello.Age = 24
	what.Hello.School = "what school2"
	fmt.Printf("what hello: %v\n", what.Hello)
	fmt.Printf("what: %v\n", what)
	fmt.Printf("what -> name: %v\n", what.Name)

	//1、Dail连接
	connect, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	defer connect.Close()
	if err != nil {
		panic(err.Error())
	}
	orderService := message.NewOrderServiceClient(connect)
	request := &message.OrderRequest{OrderId: "201907300001", TimeStamp: time.Now().Unix()}
	orderInfo, err := orderService.GetOrderInfo(context.Background(), request)
	if orderInfo != nil {
		fmt.Println(orderInfo.GetOrderId())
		fmt.Println(orderInfo.GetOrderName())
		fmt.Println(orderInfo.GetOrderStatus())
	}
}
