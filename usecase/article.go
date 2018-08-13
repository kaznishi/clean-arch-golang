package usecase

import (
	"time"

	"github.com/kaznishi/clean-arch-golang/domain/repository"
	"github.com/kaznishi/clean-arch-golang/domain/model"
)

type ArticleUsecase struct {
	articleRepository repository.ArticleRepository
	authorRepository repository.AuthorRepository
	contextTimeout time.Duration
}

func NewArticleUsecase(arRepo repository.ArticleRepository, auRepo repository.AuthorRepository, timeout time.Duration) ArticleUsecase {
	return ArticleUsecase{
		articleRepository: arRepo,
		authorRepository: auRepo,
		contextTimeout: timeout,
	}
}

type authorChannel struct {
	Author *model.Author
	Error error
}