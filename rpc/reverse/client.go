package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("ListenTcp err:", err)
	}

	clientChan := make(chan *rpc.Client)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Accept err:", err)
			}

			clientChan <- rpc.NewClient(conn)
		}
	}()

	doClientWork(clientChan)
}

func doClientWork(clientChan <-chan *rpc.Client) {
	client := <-clientChan
	defer client.Close()

	var reply string
	err := client.Call("HelloService.Hello", "rpc", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
