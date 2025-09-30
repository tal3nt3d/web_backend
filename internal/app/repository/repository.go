package repository

import (
	"errors"

	"github.com/minio/minio-go/v7"
	minioClient "web_backend/internal/app/minioClient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotAllowed    = errors.New("not allowed")
	ErrNoDraft       = errors.New("no draft for this user")
)

type Repository struct {
	db *gorm.DB
	mc     *minio.Client
	user_Id int
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mc, err := minioClient.InitMinio()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		mc: mc,
		user_Id: 0,
	}, nil
}

func (r *Repository) GetUserID() (int) {
	return r.user_Id
}

func (r *Repository) SetUserID(id int) {
	r.user_Id = id
}

func (r *Repository) SignOut() {
	r.user_Id = 0
}
