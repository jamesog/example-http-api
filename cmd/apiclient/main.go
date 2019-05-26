package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jamesog/example-http-api/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := api.NewExampleServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Request users with an empty UserRequest, i.e. all users
	// Although not implemented in the backend, this could be used to request a specific user
	users, err := client.GetUsers(ctx, &api.UserRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// For demonstration purposes, show we can still print the data as JSON in a similar format to the non-protobuf version
	j, _ := json.Marshal(users.Users)
	fmt.Printf("%s\n", j)

	for _, u := range users.Users {
		fmt.Printf("User ID %d is %s\n", u.Id, u.Name)
	}
}
