package handler

import (
	"strconv"
	"context"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
	"github.com/kaznishi/clean-arch-golang/usecase"
	"github.com/labstack/echo"
	"github.com/kaznishi/clean-arch-golang/domain/model"
)
type ResponseError struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	ArticleUsecase usecase.ArticleUsecase
}

func (a *ArticleHandler) FetchArticle(c echo.Context) error {
	numStr := c.QueryParam("num")
	num, _ := strconv.Atoi(numStr)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listArticle, nextCursor, err := a.ArticleUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listArticle)
}

func (a *ArticleHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	article, err := a.ArticleUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message:err.Error()})
	}
	return c.JSON(http.StatusOK, article)
}

func isRequestValid(a *model.Article) (bool, error) {
	validator := validator.New()
	err := validator.Struct(a)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *ArticleHandler) Store(c echo.Context) error {
	var article model.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ar, err := a.ArticleUsecase.Store(ctx, &article)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message:err.Error()})
	}
	return c.JSON(http.StatusCreated, ar)
}

func (a *ArticleHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err = a.ArticleUsecase.Delete(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message:err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case model.INTERNAL_SERVER_ERROR:
		return http.StatusInternalServerError
	case model.NOT_FOUND_ERROR:
		return http.StatusNotFound
	case model.CONFLICT_ERROR:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NewArticleHandler(e *echo.Echo, uc usecase.ArticleUsecase) {
	handler := &ArticleHandler{
		ArticleUsecase: uc,
	}

	e.GET("/article", handler.FetchArticle)
	e.POST("/article", handler.Store)
	e.GET("/article/:id", handler.GetByID)
	e.DELETE("/article/:id", handler.Delete)
}