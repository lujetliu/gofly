package main

import (
	context "context"
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
	_, err = client.Publish(
		context.Background(), &String{Value: "golang:hello go"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish(
		context.Background(), &String{Value: "docker: hello docker"})

	if err != nil {
		log.Fatal(err)
	}
}
