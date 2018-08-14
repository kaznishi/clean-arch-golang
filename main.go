package main

import (
	"fmt"
	"net/url"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/kaznishi/clean-arch-golang/infra/dao"
	"github.com/kaznishi/clean-arch-golang/usecase"
	"github.com/kaznishi/clean-arch-golang/adapter/handler"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.yaml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Tokyo")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer dbConn.Close()
	e := echo.New()
	authorRepository := dao.NewAuthorDAO(dbConn)
	articleRepository := dao.NewArticleDAO(dbConn)
	timeoutContext := time.Duration(viper.GetInt("context.timetout")) * time.Second
	articleUsecase := usecase.NewArticleUsecase(articleRepository, authorRepository, timeoutContext)
	handler.NewArticleHandler(e, articleUsecase)

	e.Start(viper.GetString("server.address"))
}
