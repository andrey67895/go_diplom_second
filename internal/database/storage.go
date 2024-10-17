package database

import (
	"context"
	"embed"
	"time"

	"github.com/google/uuid"
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

func (db DBStorageModel) GetSecretTypeByID(ctx context.Context, id string) (*model.SecretType, error) {
	var secretType model.SecretType
	err := db.DB.GetContext(ctx, &secretType, `SELECT * FROM secret_type WHERE id = $1`, id)
	return &secretType, err
}

func (db DBStorageModel) SetSecret(ctx context.Context, encoded []byte, sType string, metadata string) (model.Secret, error) {
	var secret model.Secret
	err := db.DB.GetContext(ctx, &secret, `INSERT INTO secret(encoded, type, metadata) values ($1,$2,$3) RETURNING *`, encoded, sType, metadata)
	return secret, err
}

func (db DBStorageModel) UpdateSecret(ctx context.Context, encoded []byte, sType string, metadata string, id string) (model.Secret, error) {
	var secret model.Secret
	err := db.DB.GetContext(ctx, &secret, `UPDATE secret SET encoded = $1, type = $2, metadata = $3 WHERE id = $4 RETURNING *;`, encoded, sType, metadata, id)
	return secret, err
}

func (db DBStorageModel) DeleteSecretByID(ctx context.Context, id string) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM secret  WHERE id = $1;`, id)
	return err
}

func (db DBStorageModel) InsertAuthSecretRef(ctx context.Context, authID string, secretID string) error {
	_, err := db.DB.ExecContext(ctx, `INSERT INTO auth_secret_ref values ($1,$2);`, authID, secretID)
	return err
}

func (db DBStorageModel) DeleteAuthSecretRefBySecretID(ctx context.Context, secretID string) error {
	_, err := db.DB.ExecContext(ctx, `DELETE FROM auth_secret_ref WHERE secret_id = $1;`, secretID)
	return err
}

func (db DBStorageModel) SelectSecretIDByAuthID(ctx context.Context, authID string) ([]uuid.UUID, error) {
	var secretIDs []uuid.UUID
	err := db.DB.SelectContext(ctx, &secretIDs, `select secret_id from auth_secret_ref where auth_id = $1;`, authID)
	return secretIDs, err
}

func (db DBStorageModel) GetSecretByID(ctx context.Context, ID string) (model.Secret, error) {
	var secret model.Secret
	err := db.DB.GetContext(ctx, &secret, `SELECT * FROM secret where id = $1`, ID)
	return secret, err
}
