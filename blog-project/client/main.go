package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "blogproject.com/grpc/protos"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPostsClient(conn)

	runGetAllPosts(client)
	// runGetPost(client, "1")
	// runCreatePost(client, "Learning new", "ready to learn new things daily", "aman",
	// 	"22-02-2024", []string{"Learning", "new"})
	// runUpdatePost(client, "1", "Learning new", "ready to learn new things daily", "aman",
	// 	"22-02-2024", []string{"Learning", "new"}
	// runDeletePost(client, "1")
}

func runGetAllPosts(client pb.PostsClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetAllPosts(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetAllPosts(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetAllPosts(_) = _, %v", client, err)
		}
		log.Printf("GetAllPosts: %v", row)
	}
}

func runGetPost(client pb.PostsClient, postid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: postid}
	res, err := client.GetPost(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetPost(_) = _, %v", client, err)
	}
	log.Printf("GetPost : %v", res)
}

func runCreatePost(client pb.PostsClient, title string,
	content string, author string, publicationDate string, tags []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.BlogPost{Title: title, Content: content, Author: author,
		PublicationDate: publicationDate, Tags: tags}
	res, err := client.CreatePost(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreatePost(_) = _, %v", client, err)
	}
	log.Printf("CreatePost : %v", res)
}

func runUpdatePost(client pb.PostsClient, postid string, title string,
	content string, author string, publicationDate string, tags []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.BlogPost{PostId: postid, Title: title, Content: content, Author: author,
		PublicationDate: publicationDate, Tags: tags}
	res, err := client.UpdatePost(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdatePost(_) = _, %v", client, err)
	}
	log.Printf("UpdatePost : %v", res)

}

func runDeletePost(client pb.PostsClient, postid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: postid}
	res, err := client.DeletePost(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeletePost(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeletePost Success")
	} else {
		log.Printf("DeletePost Failed")
	}
}
