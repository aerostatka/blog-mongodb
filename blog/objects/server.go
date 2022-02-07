package objects

import (
	"context"
	"fmt"
	"github.com/aerostatka/mongodb-example/blog/blogpb"
	"github.com/aerostatka/mongodb-example/blog/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blogServer struct {
	service Service
}

func (s *blogServer) CreatePost(ctx context.Context, req *blogpb.CreatePostRequest) (*blogpb.CreatePostResponse, error) {
	post := req.GetPost()
	fmt.Printf("Create post request is invoked with request: %v\n", post)
	storagePost := structs.CreatePostFromGrpcObject(*post)

	storagePost, err := s.service.CreatePost(storagePost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %s\n", err.Error()))
	}

	fmt.Printf("Assigned id for the post is: %s\n", storagePost.ID)

	return &blogpb.CreatePostResponse{
		Post: storagePost.ToBlogPbPost(),
	}, nil
}

func (s *blogServer) ReadPost(ctx context.Context, req *blogpb.ReadPostRequest) (*blogpb.ReadPostResponse, error) {
	postId := req.GetPostId()
	fmt.Printf("ReadPost post request is invoked with request: %s\n", postId)
	storagePost, err := s.service.ReadPost(postId)
	if err != nil {
		fmt.Printf("Error during object search: %s\n", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Something went wrong: %s", err.Error()),
		)
	}

	if storagePost == nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find post with id %s", postId),
		)
	}

	return &blogpb.ReadPostResponse{
		Post: storagePost.ToBlogPbPost(),
	}, nil
}

func (s *blogServer) DeletePost(ctx context.Context, req *blogpb.DeletePostRequest) (*blogpb.DeletePostResponse, error) {
	postId := req.GetPostId()
	fmt.Printf("DeletePost post request is invoked with request: %s\n", postId)
	result, err := s.service.DeletePost(postId)
	if err != nil {
		fmt.Printf("Error during object search: %s\n", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Something went wrong: %s", err.Error()),
		)
	}

	if result == false {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find post with id %s", postId),
		)
	}

	return &blogpb.DeletePostResponse{
		Result: result,
	}, nil
}

func (s *blogServer) UpdatePost(ctx context.Context, req *blogpb.UpdatePostRequest) (*blogpb.UpdatePostResponse, error) {
	post := req.GetPost()
	fmt.Printf("UpdatePost post request is invoked with request: %v\n", post)
	newContentPost := structs.CreatePostFromGrpcObject(*post)
	if newContentPost == nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot convert post object properly",
		)
	}
	updatedPost, err := s.service.UpdatePost(newContentPost)
	if err != nil {
		fmt.Printf("Error during object search: %s\n", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Something went wrong: %s", err.Error()),
		)
	}

	if updatedPost == nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find post with id %s", newContentPost.ID),
		)
	}

	return &blogpb.UpdatePostResponse{
		Post: updatedPost.ToBlogPbPost(),
	}, nil
}

func (s *blogServer) ListPosts(req *blogpb.ListPostsRequest, stream blogpb.BlogService_ListPostsServer) error {
	fmt.Println("ListPosts request is invoked")
	posts, err := s.service.ListPosts()
	if err != nil {
		fmt.Printf("Error during posts search: %s\n", err.Error())
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Something went wrong: %s", err.Error()),
		)
	}

	for _, post := range posts {
		err = stream.Send(&blogpb.ListPostsResponse{
			Post: post.ToBlogPbPost(),
		})

		if err != nil {
			fmt.Printf("Error during request send: %s\n", err.Error())
		}
	}

	return nil
}

func CreateBlogServer(s Service) *blogServer {
	return &blogServer{
		service: s,
	}
}
