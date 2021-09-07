package main

import (
	"genesis/global"
	"genesis/internal/api"
	"genesis/pkg/logging"
	"net/http"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	r := api.NewRouter()

	if global.Logged == 0 {
		r = api.NewRouter()
	} else {
		r = api.NewRouter()
	}

	logger.Info("start server")
	http.ListenAndServe(":8000", r)

}
