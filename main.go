package main

import (
	"net/http"

	"github.com/go-http-utils/logger"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello, World!!"))
	})

	http.ListenAndServe(":8080", logger.Handler(mux, os.Stdout, logger.DevLoggerType))
}