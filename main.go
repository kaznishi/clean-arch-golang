package main

import (
	"os"
	"net/http"

	"github.com/go-http-utils/logger"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello, World!!"))
	})

	http.ListenAndServe(":8080", logger.Handler(mux, os.Stdout, logger.DevLoggerType))
}