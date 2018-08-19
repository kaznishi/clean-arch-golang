package registry

import (
	"database/sql"

	"github.com/kaznishi/clean-arch-golang/domain/repository"
	"github.com/kaznishi/clean-arch-golang/infra/dao"
)

type Repository interface {
	NewArticleRepository() repository.ArticleRepository
	NewAuthorRepository() repository.AuthorRepository
}

type repositoryImpl struct {
	MySQLConn *sql.DB
}

func NewRepository(mysqlConn *sql.DB) Repository {
	return &repositoryImpl{
		MySQLConn: mysqlConn,
	}
}

func (r *repositoryImpl) NewArticleRepository() repository.ArticleRepository {
	return dao.NewArticleDAO(r.MySQLConn)
}

func (r *repositoryImpl) NewAuthorRepository() repository.AuthorRepository {
	return dao.NewAuthorDAO(r.MySQLConn)
}