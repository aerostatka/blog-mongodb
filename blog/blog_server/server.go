package main

import (
	"context"
	"fmt"
	"github.com/aerostatka/mongodb-example/blog/blogpb"
	"github.com/aerostatka/mongodb-example/blog/objects"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog server invoked!")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %s", err.Error())
	}

	fmt.Println("Connecting to MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error in MongoDB connection: %s\n", err.Error())
	}
	defer client.Disconnect(context.TODO())

	mongoRep := objects.NewMongoRepository(client)
	service := objects.NewStandardService(mongoRep)

	service.TestConnection()

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "certs/server.crt"
		keyFile := "certs/server.pem"

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed loading certificates: %s\n", err.Error())
			return
		}

		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	server := objects.CreateBlogServer(service)
	blogpb.RegisterBlogServiceServer(s, server)

	go func() {
		fmt.Println("Starting server....")
		err = s.Serve(listener)
		if err != nil {
			log.Fatalf("Failed to serve %s", err.Error())
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener...")
	listener.Close()
	fmt.Println("End of program.")
}
