package usecase

import (
	"time"

	"github.com/kaznishi/clean-arch-golang/domain/repository"
	"github.com/kaznishi/clean-arch-golang/domain/model"
	"context"
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

func (a *ArticleUsecase) getAuthorDetail(ctx context.Context, item *model.Article, authorCh chan authorChannel) {
	res, err := a.authorRepository.GetByID(ctx, item.Author.ID)
	holder := authorChannel{
		Author: res,
		Error: err,
	}
	if ctx.Err() != nil {
		return
	}
	authorCh <- holder
}


