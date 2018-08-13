package dao

import (
	"context"
	"database/sql"

	"github.com/kaznishi/clean-arch-golang/domain/repository"
	"github.com/kaznishi/clean-arch-golang/domain/model"
)

type authorDAO struct {
	DB *sql.DB
}

func NewAuthorDAO(db *sql.DB) repository.AuthorRepository {
	return &authorDAO{
		DB: db,
	}
}

func (dao *authorDAO) GetByID(ctx context.Context, id int64) (*model.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id = ?`

	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, id)
	a := &model.Author{}

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}