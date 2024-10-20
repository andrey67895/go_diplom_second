package clihelpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tcnksm/go-input"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/andrey67895/go_diplom_second/proto"
)

type KeeperCli struct {
	Log                 *zap.SugaredLogger
	Token               string
	MPassword           string
	UI                  *input.UI
	KeeperServiceClient pb.KeeperServiceClient
}

func (cli *KeeperCli) UserRegister(c *cli.Context) error {
	query := "Введите логин"
	login, _ := cli.UI.Ask(query, &input.Options{
		Required: true,
		Loop:     true,
	})
	query = "Введите пароль"
	pass, _ := cli.UI.Ask(query, &input.Options{
		Required:    true,
		Loop:        true,
		Mask:        true,
		MaskDefault: true,
	})
	query = "Введите мастер пароль"
	mPass, _ := cli.UI.Ask(query, &input.Options{
		Required:    true,
		Loop:        true,
		Mask:        true,
		MaskDefault: true,
	})

	req := pb.RegisterRequest{
		Login:          login,
		Password:       pass,
		MasterPassword: mPass,
	}
	resp, err := cli.KeeperServiceClient.Register(c.Context, &req)
	if err != nil {
		cli.Log.Infof("Ошибка регистрации: %s", err.Error())
		return err
	}
	cli.Token = resp.GetAccessToken()
	cli.MPassword = mPass
	return nil
}

func (cli *KeeperCli) UserAuth(c *cli.Context) error {
	query := "Введите логин"
	login, _ := cli.UI.Ask(query, &input.Options{
		Required: true,
		Loop:     true,
	})
	query = "Введите пароль"
	pass, _ := cli.UI.Ask(query, &input.Options{
		Required:    true,
		Loop:        true,
		Mask:        true,
		MaskDefault: true,
	})
	query = "Введите мастер пароль"
	mPass, _ := cli.UI.Ask(query, &input.Options{
		Required:    true,
		Loop:        true,
		Mask:        true,
		MaskDefault: true,
	})
	req := pb.LoginRequest{
		Login:          login,
		Password:       pass,
		MasterPassword: mPass,
	}
	resp, err := cli.KeeperServiceClient.Login(c.Context, &req)
	if err != nil {
		cli.Log.Info("Ошибка авторизации")
		return err
	}
	cli.Token = resp.GetAccessToken()
	cli.Log.Info(cli.Token)
	cli.MPassword = resp.GetAccessToken()
	return nil
}

func (cli *KeeperCli) GetAllSecret(c *cli.Context) {
	md := metadata.New(map[string]string{"token": cli.Token})
	md.Set("master_key", cli.MPassword)
	ctx := metadata.NewOutgoingContext(c.Context, md)
	allSecrets, err := cli.KeeperServiceClient.GetAllSecrets(ctx, &pb.GetAllSecretsRequest{}, grpc.Header(&md))
	if err != nil {
		cli.Log.Infof("Ошибка при получении всех секретов: %s", err.Error())
	}
	println("ID:TYPE:META")
	for _, secret := range allSecrets.Secrets {
		println(fmt.Sprintf("%s:%s:%s", secret.ID, secret.Type, secret.Meta))
	}
}

func (cli *KeeperCli) GetSecretByID(c *cli.Context) {
	query := "Введите ID:"
	id, _ := cli.UI.Ask(query, &input.Options{
		Required: true,
		Loop:     true,
	})
	req := pb.GetSecretRequest{ID: id}

	md := metadata.New(map[string]string{"token": cli.Token})
	md.Set("master_key", cli.MPassword)
	ctx := metadata.NewOutgoingContext(c.Context, md)
	secret, err := cli.KeeperServiceClient.GetSecret(ctx, &req, grpc.Header(&md))
	if err != nil {
		return
	}
	switch secret.Secret.Type {
	case "FILE":
		dir, err := os.Getwd()
		if err != nil {
			cli.Log.Error(err)
			return
		}
		println("ID:TYPE:META")
		println(fmt.Sprintf("%s:%s:%s", secret.Secret.ID, secret.Secret.Type, secret.Secret.Meta))
		filePath := fmt.Sprintf("%s/%s", dir, secret.Secret.SecretData.GetFile().GetFilename())
		create, err := os.Create(filePath)
		if err != nil {
			cli.Log.Error(err)
			return
		}
		_, err = create.Write(secret.Secret.SecretData.GetFile().GetAny())
		if err != nil {
			cli.Log.Error(err)
			return
		}
		cli.Log.Infof("Успешное сохранение файла: %s", filePath)
	case "WORD":
		println("ID:TYPE:META:TEXT")
		println(fmt.Sprintf("%s:%s:%s:%s", secret.Secret.ID, secret.Secret.Type, secret.Secret.Meta, secret.Secret.SecretData.GetText()))
	case "CARD":
		println("ID:TYPE:META:CARD_NUMBER:CARD_EXPIRED:HOLDER:CVC")
		println(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", secret.Secret.ID, secret.Secret.Type, secret.Secret.Meta, secret.Secret.SecretData.GetCreditCardData().Number, secret.Secret.SecretData.GetCreditCardData().Expired, secret.Secret.SecretData.GetCreditCardData().Holder, secret.Secret.SecretData.GetCreditCardData().CVC))
	case "LOG_PASS":
		println("ID:TYPE:META:LOGIN:PASSWORD")
		println(fmt.Sprintf("%s:%s:%s:%s:%s", secret.Secret.ID, secret.Secret.Type, secret.Secret.Meta, secret.Secret.SecretData.GetAuthentication().Login, secret.Secret.SecretData.GetAuthentication().Password))
	}
}

func (cli *KeeperCli) SaveSecret(c *cli.Context) {
	query := "Тип данных для сохранения:"
	step, _ := cli.UI.Select(query, []string{"Банковская карта", "Текст", "Файл", "Логин/Пароль", "Назад"}, &input.Options{
		Loop: true,
	})
	switch step {
	case "Банковская карта":
		query := "Введите CardNumber:"
		number, _ := cli.UI.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
		})
		query = "Введите CardExpired(формат - 01/11)(Поле необязательное)"
		expired, _ := cli.UI.Ask(query, &input.Options{
			Required: false,
			Loop:     true,
		})
		query = "Введите Holder (Поле необязательное)"
		holder, _ := cli.UI.Ask(query, &input.Options{
			Required: false,
			Loop:     true,
		})
		query = "Введите CVC (Поле необязательное)"
		cvc, _ := cli.UI.Ask(query, &input.Options{
			Required:    false,
			Loop:        true,
			Mask:        true,
			MaskDefault: true,
		})
		req := pb.SetSecretRequest{Secret: &pb.Secret{
			Type: "CARD",
			Meta: "",
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_CreditCardData{CreditCardData: &pb.CreditCardData{
				Number:  number,
				Expired: expired,
				Holder:  holder,
				CVC:     cvc,
			}}},
		}}
		md := metadata.New(map[string]string{"token": cli.Token})
		md.Set("master_key", cli.MPassword)
		ctx := metadata.NewOutgoingContext(c.Context, md)
		_, err := cli.KeeperServiceClient.SetSecret(ctx, &req, grpc.Header(&md))
		if err != nil {
			cli.Log.Error(err)
		}
		return
	case "Текст":
		query := "Введите текст:"
		text, _ := cli.UI.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
		})
		req := pb.SetSecretRequest{Secret: &pb.Secret{
			Type:       "WORD",
			Meta:       "",
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_Text{Text: text}},
		}}
		md := metadata.New(map[string]string{"token": cli.Token})
		md.Set("master_key", cli.MPassword)
		ctx := metadata.NewOutgoingContext(c.Context, md)
		_, err := cli.KeeperServiceClient.SetSecret(ctx, &req, grpc.Header(&md))
		if err != nil {
			cli.Log.Error(err)
		}
		return
	case "Файл":
		query := "Введите путь до файла:"
		filePath, _ := cli.UI.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
		})

		file, err := os.ReadFile(filePath)
		if err != nil {
			cli.Log.Infof("Ошибка при чтении файла: %s", err.Error())
		}
		req := pb.SetSecretRequest{Secret: &pb.Secret{
			Type: "FILE",
			Meta: "",
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_File{File: &pb.File{
				Filename: filepath.Base(filePath),
				Any:      file,
			}}},
		}}
		md := metadata.New(map[string]string{"token": cli.Token})
		md.Set("master_key", cli.MPassword)
		ctx := metadata.NewOutgoingContext(c.Context, md)
		_, err = cli.KeeperServiceClient.SetSecret(ctx, &req, grpc.Header(&md))
		if err != nil {
			cli.Log.Error(err)
		}
		return
	case "Логин/Пароль":
		query := "Введите login:"
		login, _ := cli.UI.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
		})
		query = "Введите password"
		pass, _ := cli.UI.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
		})
		req := pb.SetSecretRequest{Secret: &pb.Secret{
			Type: "LOG_PASS",
			Meta: "",
			SecretData: &pb.Secret_Data{Variant: &pb.Secret_Data_Authentication{Authentication: &pb.AuthenticationData{
				Login:    login,
				Password: pass,
			}}},
		}}
		md := metadata.New(map[string]string{"token": cli.Token})
		md.Set("master_key", cli.MPassword)
		ctx := metadata.NewOutgoingContext(c.Context, md)
		_, err := cli.KeeperServiceClient.SetSecret(ctx, &req, grpc.Header(&md))
		if err != nil {
			cli.Log.Error(err)
		}
		return
	case "Назад":
		return
	}
}

func (cli *KeeperCli) RemoveSecretByID(c *cli.Context) {
	query := "Введите ID:"
	id, _ := cli.UI.Ask(query, &input.Options{
		Required: true,
		Loop:     true,
	})
	req := pb.RemoveSecretRequest{ID: id}
	md := metadata.New(map[string]string{"token": cli.Token})
	ctx := metadata.NewOutgoingContext(c.Context, md)
	_, err := cli.KeeperServiceClient.RemoveSecret(ctx, &req, grpc.Header(&md))
	if err != nil {
		cli.Log.Error(err.Error())
		return
	}
}
