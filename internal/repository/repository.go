package repository

import (
	db "github.com/0utl1er-tech/prism-backend/gen/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IRepository interface {
	db.Querier
}

type Repository struct {
	q db.Querier
}

func NewRepository(connPool *pgxpool.Pool) IRepository {
	return &Repository{q: db.New(connPool)}
}
