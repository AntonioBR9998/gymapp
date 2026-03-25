package core

import (
	"github.com/AntonioBR9998/gymapp/internal/config"
	"github.com/AntonioBR9998/gymapp/internal/repository"
)

type Core interface {
}

type core struct {
	repo repository.Repository
	conf config.Config
}

func NewCore(repo repository.Repository, conf config.Config) Core {
	c := &core{
		repo: repo,
		conf: conf,
	}

	return c
}
