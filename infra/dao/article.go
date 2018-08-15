package dao

import (
	"database/sql"
	"context"
	"fmt"

	"github.com/kaznishi/clean-arch-golang/domain/model"
	"github.com/kaznishi/clean-arch-golang/domain/repository"
)

type articleDAO struct {
	DB *sql.DB
}

func NewArticleDAO(DB *sql.DB) repository.ArticleRepository {
	return &articleDAO{DB}
}

func (dao *articleDAO) fetch(ctx context.Context, query string, args ...interface{}) ([]*model.Article, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*model.Article, 0)
	for rows.Next() {
		t := new(model.Article)
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		t.Author = model.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}


func (dao *articleDAO) Fetch(ctx context.Context, cursor string, num int64) ([]*model.Article, error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE id > ? LIMIT ?`

	return dao.fetch(ctx, query, cursor, num)
}

func (dao *articleDAO) GetByID(ctx context.Context, id int64) (*model.Article, error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE id = ?`

	list, err := dao.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	a := &model.Article{}
	if len(list) > 0 {
		a = list[0]
	} else {
		return nil, model.NOT_FOUND_ERROR
	}

	return a, nil
}

func (dao *articleDAO) GetByTitle(ctx context.Context, title string) (*model.Article, error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE title = ?`

	list, err := dao.fetch(ctx, query, title)
	if err != nil {
		return nil, err
	}

	a := &model.Article{}
	if  len(list) > 0 {
		a = list[0]
	} else {
		return nil, model.NOT_FOUND_ERROR
	}

	return a, nil
}

func (dao *articleDAO) Store(ctx context.Context, a *model.Article) (int64, error) {
	query := `INSERT article SET title = ?, content = ?, author_id = ?, updated_at = ?, created_at = ?`
	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (dao *articleDAO) Delete(ctx context.Context, id int64) (bool, error) {
	query := `DELETE FROM article WHERE id = ?`

	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", rowsAffected)
		return false, err
	}

	return true, nil
}

func (dao *articleDAO) Update(ctx context.Context, a *model.Article) (*model.Article, error) {
	query := `UPDATE article SET title = ?, content = ?, author_id = ?, updated_at = ? WHERE id = ?`

	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, nil
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.ID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", rowsAffected)
		return nil, err
	}

	return a, nil
}





