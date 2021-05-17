package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var sb = NewBuilder()

// Repository represents type to handle database operations
type Repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates new instance of Repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// NewBuilder returns new pgSQL query builder
func NewBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
