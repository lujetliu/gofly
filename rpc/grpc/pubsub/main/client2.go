package main

import (
	context "context"
	"fmt"
	"io"
	"log"

	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewPubsubServiceClient(conn)
	stream, err := client.Subscribe(
		context.Background(), &String{Value: "golang:"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
}
