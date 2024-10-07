package services

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andrey67895/go_diplom_second/internal/database"
	"github.com/andrey67895/go_diplom_second/internal/helpers"
	"github.com/andrey67895/go_diplom_second/internal/logger"
	pb "github.com/andrey67895/go_diplom_second/proto"
)

var log = logger.Logger()

type KeeperServer struct {
	pb.UnimplementedKeeperServiceServer
	DB database.DBStorageModel
}

func (ks *KeeperServer) GetPing(context.Context, *pb.Ping) (*pb.Ping, error) {
	err := ks.DB.DB.Ping()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "failed to get ping")
	}
	return &pb.Ping{}, nil
}

func (ks *KeeperServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := ks.DB.Register(ctx, req.GetLogin(), req.GetPassword(), req.GetMasterPassword())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unavailable, "ошибка при регистрации")
	}
	token, err := helpers.GenerateJWTAndCheck(req.Login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unavailable, "ошибка при формировании юзера")
	}
	return &pb.RegisterResponse{AccessToken: token}, nil
}

func (ks *KeeperServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	authData, err := ks.DB.GetAuthData(ctx, req.Login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unauthenticated, "ошибка авторизации")
	}
	if authData.HashPass == helpers.EncodeHashSha256(req.GetPassword()) {
		token, err := helpers.GenerateJWTAndCheck(req.Login)
		if err != nil {
			log.Error(err)
			return nil, status.Error(codes.Unavailable, "ошибка при формировании юзера")
		}
		return &pb.LoginResponse{AccessToken: token}, nil
	} else {
		return nil, status.Error(codes.Unauthenticated, "ошибка авторизации")
	}
}

func (ks *KeeperServer) SetSecret(ctx context.Context, req *pb.SetSecretRequest) (*pb.SetSecretResponse, error) {
	_, err := ks.DB.GetSecretType(ctx, req.GetSecret().GetType())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "неверный тип секрета")
	}
	//login := "qwer"
	//authData, err := ks.DB.GetAuthData(ctx, login)
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Unavailable, "ошибка данных")
	//}
	//helpers.EncodeHashSha512(authData.HashPassmaster)

	aes, err := helpers.EncryptAES([]byte("qwerrty"), helpers.EncodeHashSha512("q2312312312312"))
	log.Error(err)
	log.Info(string(aes))
	aes, err = helpers.DecryptAES(aes, helpers.EncodeHashSha512("q2312312312312"))
	log.Info(string(aes))
	//switch expr {
	//
	//}
	//ks.DB.SetSecret(ctx)
	return nil, nil
}
