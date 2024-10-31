package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/andrey67895/go_diplom_second/internal/config"
	"github.com/andrey67895/go_diplom_second/internal/database"
	"github.com/andrey67895/go_diplom_second/internal/logger"
	"github.com/andrey67895/go_diplom_second/internal/services"
	pb "github.com/andrey67895/go_diplom_second/proto"
)

var log = logger.Logger()

func InitServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	config.InitServerConfig()
	listen, err := net.Listen("tcp", config.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	var st *database.DBStorageModel
	if config.DatabaseDsn != "" {
		st, err = database.InitDB(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	pb.RegisterKeeperServiceServer(s, &services.KeeperServer{
		DB: *st,
	})
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatal("listen and serve returned err: %v", err)
		}
	}()
	<-ctx.Done()
	log.Info("got interruption signal")
	s.GracefulStop()
	log.Info("final")
}
