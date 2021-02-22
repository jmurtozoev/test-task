package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/jmurtozoev/test-task/storage/postgres"
	"github.com/jmurtozoev/test-task/storage/repo"
)

type Storage interface {
	Product() repo.Product
}

func New(db *sqlx.DB) Storage {
	return &storage{
		productRepo: postgres.NewProductRepo(db),
	}
}

type storage struct {
	productRepo repo.Product
}

func (s *storage) Product() repo.Product {
	return s.productRepo
}