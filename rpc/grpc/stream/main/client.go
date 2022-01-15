package main

import (
	context "context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &String{Value: "grpc"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 将发送和接收操作放到两个独立的 goroutine
	go func() {
		for {
			if err := stream.Send(&String{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // TODO: io.EOF 原理, 如何模拟
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}

}
