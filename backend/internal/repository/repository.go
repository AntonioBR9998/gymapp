package repository

import "github.com/AntonioBR9998/gymapp/internal/config"

type Repository interface {
}

type repository struct {
}

func NewRepository(cfg config.Config) Repository {
	return &repository{}
}
