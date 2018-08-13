package repository

import (
	"context"

	"github.com/kaznishi/clean-arch-golang/domain/model"
)

type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*model.Article, error)
	GetByID(ctx context.Context, id int64) (*model.Article, error)
	GetByTitle(ctx context.Context, title string) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) (*model.Article, error)
	Store(ctx context.Context, a *model.Article) (int64, error)
	Delete(ctx context.Context, id int64) (bool, error)
}