package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "blogproject.com/grpc/protos"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var blogPosts []*pb.BlogPost

type blogPostServer struct {
	pb.UnimplementedPostsServer
}

func main() {
	initPosts()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterPostsServer(s, &blogPostServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initPosts() {
	post1 := &pb.BlogPost{PostId: "1", Title: "Learn grpc", Content: "Learn grpc with me", Author: "John",
		PublicationDate: "24-03-2024", Tags: []string{"golang", "grpc"},
	}
	post2 := &pb.BlogPost{PostId: "2", Title: "Learn rest", Content: "Learn rest with me", Author: "Kevin",
		PublicationDate: "20-03-2024", Tags: []string{"golang", "rest", "apis"},
	}

	blogPosts = append(blogPosts, post1, post2)
}

func (s *blogPostServer) GetAllPosts(in *pb.Empty,
	stream pb.Posts_GetAllPostsServer) error {
	log.Printf("Received: %v", in)
	for _, post := range blogPosts {
		if err := stream.Send(post); err != nil {
			return err
		}
	}
	return nil
}

func (s *blogPostServer) GetPost(ctx context.Context,
	in *pb.Id) (*pb.BlogPost, error) {
	log.Printf("Received: %v", in)

	res := &pb.BlogPost{}

	for _, post := range blogPosts {
		if post.GetPostId() == in.GetValue() {
			res = post
			break
		}
	}

	return res, nil
}

func (s *blogPostServer) CreatePost(ctx context.Context,
	in *pb.BlogPost) (*pb.BlogPost, error) {
	log.Printf("Received: %v", in)
	res := &pb.BlogPost{}
	res.PostId = strconv.Itoa(rand.Intn(100000000))
	in.PostId = res.GetPostId()
	blogPosts = append(blogPosts, in)
	res = in
	return res, nil
}

func (s *blogPostServer) UpdatePost(ctx context.Context,
	in *pb.BlogPost) (*pb.BlogPost, error) {
	log.Printf("Received: %v", in)

	res := &pb.BlogPost{}
	for index, post := range blogPosts {
		if post.GetPostId() == in.GetPostId() {
			blogPosts = append(blogPosts[:index], blogPosts[index+1:]...)
			in.PostId = post.GetPostId()
			blogPosts = append(blogPosts, in)
			res = in
			break
		}
	}

	return res, nil
}

func (s *blogPostServer) DeletePost(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}
	for index, post := range blogPosts {
		if post.GetPostId() == in.GetValue() {
			blogPosts = append(blogPosts[:index], blogPosts[index+1:]...)
			res.Value = 1
			break
		}
	}

	return &res, nil
}
