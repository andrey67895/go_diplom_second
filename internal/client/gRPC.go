package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/andrey67895/go_diplom_second/proto"
)

func NewClient(addr string) pb.KeeperServiceClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Отсутсвует доступ до сервера")
	}
	return pb.NewKeeperServiceClient(conn)
}
