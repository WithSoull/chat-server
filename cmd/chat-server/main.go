package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/withsoull/chat-server-grpc/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Println("Chat was created for:\t", req.GetUsernames())
	return &desc.CreateResponse{
		Id: int64(gofakeit.Number(1, 100)),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("Chat(ID: %d) has been deleted", req.GetId())
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	fmt.Printf("Message at '%s' from '%s':\n%s", req.GetSendAt().AsTime().Format(time.RFC3339), req.GetFrom(), req.GetText())
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
