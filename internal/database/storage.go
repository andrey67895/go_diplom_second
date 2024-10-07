package database

import (
	"context"
	"embed"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/andrey67895/go_diplom_second/internal/config"
	"github.com/andrey67895/go_diplom_second/internal/database/migrator"
	"github.com/andrey67895/go_diplom_second/internal/helpers"
	"github.com/andrey67895/go_diplom_second/internal/model"
)

type DBStorageModel struct {
	DB  *sqlx.DB
	ctx context.Context
}

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const migrationsDir = "migrations"

func openDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", config.DatabaseDsn)
	if err != nil {

	}
	return db, err
}

func InitDB(ctx context.Context) (*DBStorageModel, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(tCtx); err != nil {
		return nil, err
	}
	dbStorage := DBStorageModel{DB: db, ctx: ctx}

	tMigrator := migrator.NewMigrator(MigrationsFS, migrationsDir)
	db, err = openDB()
	if err != nil {
		return nil, err
	}
	err = tMigrator.ApplyMigrations(db)
	if err != nil {
		return nil, err
	}
	return &dbStorage, nil
}

func (db DBStorageModel) Register(ctx context.Context, login string, pass string, masterPass string) error {
	_, err := db.DB.ExecContext(ctx, `INSERT INTO auth(login, hash_pass, hash_pass_master) values ($1,$2, $3)`, login, helpers.EncodeHashSha256(pass), helpers.EncodeHashSha256(masterPass))
	return err
}

func (db DBStorageModel) GetAuthData(ctx context.Context, login string) (*model.AuthData, error) {
	var authData model.AuthData
	err := db.DB.GetContext(ctx, &authData, `SELECT * FROM auth WHERE login = $1`, login)
	return &authData, err
}

func (db DBStorageModel) GetSecretType(ctx context.Context, name string) (*model.SecretType, error) {
	var secretType model.SecretType
	err := db.DB.GetContext(ctx, &secretType, `SELECT * FROM secret_type WHERE name = $1`, name)
	return &secretType, err
}

func (db DBStorageModel) SetSecret(ctx context.Context, encoded []byte, sType string) error {
	_, err := db.DB.ExecContext(ctx, `INSERT INTO secret(encoded, type) values ($1,$2)`, encoded, sType)
	return err
}
