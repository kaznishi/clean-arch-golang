package repository

import (
	"context"

	"github.com/kaznishi/clean-arch-golang/domain/model"
)

type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Author, error)
}
