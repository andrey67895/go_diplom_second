package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"slices"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/andrey67895/go_diplom_second/internal/database"
	"github.com/andrey67895/go_diplom_second/internal/helpers"
	"github.com/andrey67895/go_diplom_second/internal/logger"
	"github.com/andrey67895/go_diplom_second/internal/model"
	"github.com/andrey67895/go_diplom_second/internal/model/data_type_model"
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
		return nil, status.Error(codes.Unauthenticated, "ошибка при регистрации")
	}
	token, err := helpers.GenerateJWTAndCheck(req.Login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unauthenticated, "ошибка при формировании токена")
	}
	return &pb.RegisterResponse{AccessToken: token}, nil
}

func (ks *KeeperServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	authData, err := ks.DB.GetAuthData(ctx, req.Login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unauthenticated, "ошибка авторизации")
	}
	if authData.HashPass == helpers.EncodeHashSha256(req.GetPassword()) && authData.HashPassmaster == helpers.EncodeHashSha256(req.GetMasterPassword()) {
		token, err := helpers.GenerateJWTAndCheck(req.Login)
		if err != nil {
			log.Error(err)
			return nil, status.Error(codes.Unauthenticated, "ошибка при формировании токена")
		}
		return &pb.LoginResponse{AccessToken: token}, nil
	} else {
		return nil, status.Error(codes.Unauthenticated, "ошибка авторизации")
	}
}

func (ks *KeeperServer) SetSecret(ctx context.Context, req *pb.SetSecretRequest) (*pb.SetSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "failed to get metadata")
	}
	login, err := getToken(md)
	if err != nil {
		return nil, err
	}
	masterKey := md["master_key"]
	if len(masterKey) == 0 {

	}
	secretType, err := ks.DB.GetSecretType(ctx, req.GetSecret().GetType())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "неверный тип секрета")
	}
	marshal, err := json.Marshal(req.Secret.SecretData.Variant)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "неверный тип секрета")
	}
	authData, err := ks.DB.GetAuthData(ctx, login)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "ошибка получения информации об юзере")
	}
	encoded := base64.StdEncoding.EncodeToString(helpers.Compress(marshal))
	byEncoded := helpers.EncryptDecrypt(encoded, helpers.EncodeHashSha512(masterKey[0]))
	secret, err := ks.DB.SetSecret(ctx, []byte(byEncoded), secretType.ID.String(), req.GetSecret().GetMeta())
	if err != nil {
		return nil, status.Error(codes.Unavailable, "ошибка сохранения секрета")
	}
	err = ks.DB.InsertAuthSecretRef(ctx, authData.ID.String(), secret.ID.String())
	if err != nil {
		return nil, status.Error(codes.Unavailable, "ошибка сохранения связи юзера и секрета")
	}
	return &pb.SetSecretResponse{}, nil
}

func (ks *KeeperServer) GetSecret(ctx context.Context, req *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "failed to get metadata")
	}
	login, err := getToken(md)
	if err != nil {
		return nil, err
	}
	masterKeys := md["master_key"]
	var masterKey string
	if len(masterKeys) != 0 {
		masterKey = masterKeys[0]
	}
	authData, err := ks.DB.GetAuthData(ctx, login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unavailable, "ошибка получения информации об юзере")
	}
	secretIDs, err := ks.DB.SelectSecretIDByAuthID(ctx, authData.ID.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ошибка сохранения секрета")
	}
	getID, err := uuid.Parse(req.GetID())
	if err != nil || !slices.Contains(secretIDs, getID) {
		return nil, status.Error(codes.NotFound, "Отсутсвтвуют сохранненные секреты")
	}
	secret, err := ks.DB.GetSecretByID(ctx, req.GetID())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.NotFound, "Отсутсвтвуют сохранненные секреты")
	}
	secretType, err := ks.DB.GetSecretTypeByID(ctx, secret.Type.String())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "оОшибка при получении типа секрета")
	}
	byDecoded := helpers.EncryptDecrypt(string(secret.Encoded), helpers.EncodeHashSha512(masterKey))
	decoded, _ := base64.StdEncoding.DecodeString(byDecoded)
	decompress, _ := helpers.DeCompress(decoded)
	return &pb.GetSecretResponse{Secret: getSecretDataByType(secretType.Name, secret, decompress)}, nil
}

func (ks *KeeperServer) GetAllSecrets(ctx context.Context, _ *pb.GetAllSecretsRequest) (*pb.GetAllSecretsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "failed to get metadata")
	}
	login, err := getToken(md)
	if err != nil {
		return nil, err
	}
	masterKeys := md["master_key"]
	var masterKey string
	if len(masterKeys) != 0 {
		masterKey = masterKeys[0]
	}
	authData, err := ks.DB.GetAuthData(ctx, login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unavailable, "ошибка получения информации об юзере")
	}
	secretIDs, err := ks.DB.SelectSecretIDByAuthID(ctx, authData.ID.String())
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, "ошибка при получении связи юзера и секрета")
	}
	var secrets []*pb.Secret
	for _, id := range secretIDs {
		secret, err := ks.DB.GetSecretByID(ctx, id.String())
		if err != nil {
			log.Error(err)
			return nil, status.Error(codes.InvalidArgument, "ошибка при получении секрета")
		}
		secretType, err := ks.DB.GetSecretTypeByID(ctx, secret.Type.String())
		if err != nil {
			log.Error(err)
			return nil, status.Error(codes.InvalidArgument, "неверный тип секрета")
		}
		byDecoded := helpers.EncryptDecrypt(string(secret.Encoded), helpers.EncodeHashSha512(masterKey))
		decoded, _ := base64.StdEncoding.DecodeString(byDecoded)
		decompress, _ := helpers.DeCompress(decoded)
		secrets = append(secrets, getSecretDataByType(secretType.Name, secret, decompress))
	}
	response := pb.GetAllSecretsResponse{}
	response.Secrets = secrets
	return &response, nil
}

func (ks *KeeperServer) RemoveSecret(ctx context.Context, req *pb.RemoveSecretRequest) (*pb.RemoveSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "Ошибка при получении токена")
	}
	login, err := getToken(md)
	if err != nil {
		return nil, err
	}
	authData, err := ks.DB.GetAuthData(ctx, login)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Unavailable, "ошибка получения информации об юзере")
	}
	secretIDs, err := ks.DB.SelectSecretIDByAuthID(ctx, authData.ID.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ошибка при получении связи юзера и секрета")
	}
	deletedID, err := uuid.Parse(req.GetID())
	if err != nil || !slices.Contains(secretIDs, deletedID) {
		return nil, status.Error(codes.InvalidArgument, "ошибка удаления секрета")
	}
	err = ks.DB.DeleteAuthSecretRefBySecretID(ctx, deletedID.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ошибка удаления секрета")
	}
	err = ks.DB.DeleteSecretByID(ctx, req.GetID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ошибка удаления секрета")
	}
	return &pb.RemoveSecretResponse{}, nil
}

func getToken(md metadata.MD) (string, error) {
	token := md["token"]
	login := ""
	if len(token) != 0 {
		var err error
		login, err = helpers.DecodeJWT(token[0])
		if err != nil {
			return "", status.Error(codes.Unauthenticated, "ошибка авторизации")
		}
	}
	if login == "" {
		return "", status.Error(codes.Unauthenticated, "ошибка авторизации")
	}
	return login, nil
}

func getSecretDataByType(tp string, baseSecret model.Secret, model []byte) *pb.Secret {
	switch tp {
	case "LOG_PASS":
		var login data_type_model.AuthData
		json.Unmarshal(model, &login)
		return &pb.Secret{
			ID:   baseSecret.ID.String(),
			Type: tp,
			Meta: baseSecret.Metadata,
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_Authentication{Authentication: &pb.AuthenticationData{
				Login:    login.Authentication.Login,
				Password: login.Authentication.Password,
			}}}}
	case "FILE":
		var file data_type_model.FileData
		json.Unmarshal(model, &file)
		return &pb.Secret{
			ID:   baseSecret.ID.String(),
			Type: tp,
			Meta: baseSecret.Metadata,
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_File{File: &pb.File{
				Any:      file.File.Data,
				Filename: file.File.Filename,
			}}}}
	case "WORD":
		var word data_type_model.Word
		json.Unmarshal(model, &word)
		return &pb.Secret{
			ID:         baseSecret.ID.String(),
			Type:       tp,
			Meta:       baseSecret.Metadata,
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_Text{Text: word.Text}}}
	case "CARD":
		var card data_type_model.CardData

		json.Unmarshal(model, &card)
		return &pb.Secret{
			ID:   baseSecret.ID.String(),
			Type: tp,
			Meta: baseSecret.Metadata,
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_CreditCardData{CreditCardData: &pb.CreditCardData{
				Number:  card.Card.Number,
				Expired: card.Card.Expired,
				Holder:  card.Card.Holder,
				CVC:     card.Card.CVC,
			}}}}
	}
	return nil
}
