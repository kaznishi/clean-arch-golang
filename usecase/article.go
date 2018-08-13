package usecase

import (
	"time"
	"context"
	"strconv"

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

func (a *ArticleUsecase) getAuthorDetails(ctx context.Context, data []*model.Article) ([]*model.Article, error) {
	authorCh := make(chan authorChannel)
	defer close(authorCh)
	existingAuthorMap := make(map[int64]bool)
	for _, item := range data {
		if _, ok := existingAuthorMap[item.Author.ID]; !ok {
			existingAuthorMap[item.Author.ID] = true
			go a.getAuthorDetail(ctx, item, authorCh)
		}
	}

	mapAuthor := make(map[int64]*model.Author)
	totalGoroutineCalled := len(existingAuthorMap)
	for i := 0; i < totalGoroutineCalled; i++ {
		select {
		case a := <-authorCh:
			if a.Error == nil {
				if a.Author != nil {
					mapAuthor[a.Author.ID] = a.Author
				}
			} else {
				return nil, a.Error
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	for index, item := range data {
		if a, ok := mapAuthor[item.Author.ID]; ok {
			data[index].Author = *a
		}
	}

	return data, nil
}

func (a *ArticleUsecase) Fetch(c context.Context, cursor string, num int64) ([]*model.Article, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listArticle, err := a.articleRepository.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""

	listArticle, err = a.getAuthorDetails(ctx, listArticle)
	if err != nil {
		return nil, "", err
	}

	if size := len(listArticle); size == int(num) {
		lastId := listArticle[num - 1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	return listArticle, nextCursor, nil
}

func (a *ArticleUsecase) GetByID(c context.Context, id int64) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepository.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor
	return res, nil
}

func (a *ArticleUsecase) Update(c context.Context, ar *model.Article) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepository.Update(ctx, ar)
}

func (a *ArticleUsecase) GetByTitle(c context.Context, title string) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.articleRepository.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepository.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor

	return res, nil
}

func (a *ArticleUsecase) Store(c context.Context, ar *model.Article) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existArticle, _ := a.GetByTitle(ctx, ar.Title)
	if existArticle != nil {
		return nil, model.CONFLICT_ERROR
	}

	id, err := a.articleRepository.Store(ctx, ar)
	if err != nil {
		return nil, err
	}
	ar.ID = id
	return ar, nil
}

func (a *ArticleUsecase) Delete(c context.Context, id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existArticle, _ := a.articleRepository.GetByID(ctx, id)
	if existArticle == nil {
		return false, model.NOT_FOUND_ERROR
	}
	return a.articleRepository.Delete(ctx, id)
}
