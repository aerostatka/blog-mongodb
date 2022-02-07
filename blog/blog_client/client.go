package main

import (
	"context"
	"fmt"
	"github.com/aerostatka/mongodb-example/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func doCreatePost(c blogpb.BlogServiceClient) {
	fmt.Println("doCreatePost is running....")

	post := &blogpb.Post{
		AuthorId: "Daria",
		Title:    "My first post",
		Content:  "No post content is available.",
	}

	response, err := c.CreatePost(context.Background(), &blogpb.CreatePostRequest{
		Post: post,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %s\n", err.Error())
	}

	fmt.Printf("Received ID of the new post: %s\n", response.GetPost().GetId())
}

func doReadPost(c blogpb.BlogServiceClient) {
	fmt.Println("doReadPost is running....")

	requests := []*blogpb.ReadPostRequest{
		{
			PostId: "12",
		},
		{
			PostId: "61f7452dbf246465380baf54",
		},
		{
			PostId: "61f7452dbf246465380baf5d",
		},
	}

	for _, request := range requests {
		response, err := c.ReadPost(context.Background(), request)
		if err != nil {
			fmt.Printf("Unexpected error: %s\n", err.Error())
			continue
		}

		fmt.Printf("Received post information: %v\n", response.GetPost())
	}
}

func doUpdatePost(c blogpb.BlogServiceClient) {
	fmt.Println("doUpdatePost is running....")

	post := &blogpb.Post{
		Id:       "61fb53e09fde6738200c0ca1",
		AuthorId: "Alexander",
		Title:    "My updated post",
		Content:  "Content is now available.",
	}

	response, err := c.UpdatePost(context.Background(), &blogpb.UpdatePostRequest{
		Post: post,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %s\n", err.Error())
	}

	fmt.Printf("Received ID of the updated post: %s\n", response.GetPost().GetId())
}

func doDeletePost(c blogpb.BlogServiceClient) {
	fmt.Println("doDeletePost is running....")

	requests := []*blogpb.DeletePostRequest{
		{
			PostId: "61f7459ebf246465380baf5e",
		},
		{
			PostId: "61f7452dbf246465380baf5d",
		},
	}

	for _, request := range requests {
		response, err := c.DeletePost(context.Background(), request)
		if err != nil {
			fmt.Printf("Unexpected error: %s\n", err.Error())
			continue
		}

		fmt.Printf("Status of post deletion: %v\n", response.GetResult())
	}
}

func doListPosts(c blogpb.BlogServiceClient) {
	fmt.Println("doListPosts is running....")

	stream, err := c.ListPosts(context.Background(), &blogpb.ListPostsRequest{})
	if err != nil {
		log.Fatalf("Error during sending lest request: %s\n", err.Error())
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error during receiving the request: %s\n", err.Error())
		}

		fmt.Printf("Received post content: %v\n", result.GetPost())
	}
}

func main() {
	fmt.Println("Client has been invoked!")
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	tls := false
	if tls {
		creds, err := credentials.NewClientTLSFromFile("certs/ca.crt", "")
		if err != nil {
			log.Fatalf("Error while loading certificates: %v", err.Error())
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err.Error())
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	doCreatePost(c)
	doReadPost(c)
	doUpdatePost(c)
	doDeletePost(c)
	doListPosts(c)
}
