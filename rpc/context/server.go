package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go func() {
			defer conn.Close()
			h := rpc.NewServer()
			h.Register(&HelloService{conn: conn})
			h.ServeConn(conn)
		}()
	}
}
