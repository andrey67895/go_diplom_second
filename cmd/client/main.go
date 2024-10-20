package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tcnksm/go-input"
	"github.com/urfave/cli/v2"

	"github.com/andrey67895/go_diplom_second/internal/client"
	"github.com/andrey67895/go_diplom_second/internal/clihelpers"
	"github.com/andrey67895/go_diplom_second/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	defer stop()
	log := logger.Logger()

	go func() {
		<-ctx.Done()
		os.Exit(0)
	}()
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server_addr",
				Aliases: []string{"sa"},
			},
			&cli.StringFlag{
				Name:    "login",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "mPassword",
				Aliases: []string{"mp"},
			},
			&cli.StringFlag{
				Name: "secretID",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "init",
				Action: func(c *cli.Context) error {
					ui := &input.UI{
						Writer: os.Stdout,
						Reader: os.Stdin,
					}
					keeperServiceClient := client.NewClient(c.String("server_addr"))
					keeperCli := clihelpers.KeeperCli{
						Log:                 log,
						Token:               "",
						MPassword:           "",
						UI:                  ui,
						KeeperServiceClient: keeperServiceClient,
					}
					for {
						if keeperCli.Token == "" {
							query := "Выберите действие:"
							step, _ := ui.Select(query, []string{"Зарегистрироваться", "Авторизоваться", "Выйти"}, &input.Options{
								Loop: true,
							})
							switch step {
							case "Зарегистрироваться":
								err := keeperCli.UserRegister(c)
								if err != nil {
									continue
								}
							case "Авторизоваться":
								err := keeperCli.UserAuth(c)
								if err != nil {
									continue
								}
							case "Выйти":
								log.Info("Завершение работы GophKeeper")
								return nil
							}
						}
						for {
							query := "Выберите действие:"
							step, err := ui.Select(query, []string{"Получить список данных", "Получить данные по ID", "Сохранить данные", "Удалить данные по ID", "Сменить пользователя", "Выйти"}, &input.Options{
								Loop: true,
							})
							if err != nil {
								log.Error(err)
								continue
							}
							switch step {
							case "Получить список данных":
								keeperCli.GetAllSecret(c)
								continue
							case "Получить данные по ID":
								keeperCli.GetSecretByID(c)
								continue
							case "Сохранить данные":
								keeperCli.SaveSecret(c)
								continue
							case "Удалить данные по ID":
								keeperCli.RemoveSecretByID(c)
								continue
							case "Сменить пользователя":
								keeperCli.Token = ""
								keeperCli.MPassword = ""
							case "Выйти":
								log.Info("Завершение работы GophKeeper")
								return nil
							}
							if keeperCli.Token == "" {
								break
							}
						}
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
