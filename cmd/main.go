package main

import (
	"log"
	"net/http"
	"startdownrec"

	"github.com/unrolled/logger"
)

func main() {
	loggerMiddleware := logger.New(logger.Options{
		Prefix: "StartDownRec",
		RemoteAddressHeaders: []string{"X-Forwarded-For"},
		OutputFlags: log.LstdFlags,
	})

	mux := http.NewServeMux()
	mux.Handle("/", loggerMiddleware.Handler(http.HandlerFunc(startdownrec.Run)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
