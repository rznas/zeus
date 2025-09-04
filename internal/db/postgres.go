package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOptions struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

func NewPostgres(opts PostgresOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", opts.Host, opts.User, opts.Password, opts.DB, opts.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
