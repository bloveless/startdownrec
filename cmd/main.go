package main

import (
	"database/sql"
	"log"
	"net/http"
	"startdownrec"

	"github.com/unrolled/logger"
)

func main() {
	dbConn, err := sql.Open("mysql", "startdownrec:2i2XfaLV4u!P1nadznz@tcp(db:3306)/startdownrec")
	if err != nil {
		log.Fatal("Unable to connect to mysql")
	}
	defer dbConn.Close()

	f := startdownrec.Function{
		DBConn: dbConn,
	}

	loggerMiddleware := logger.New(logger.Options{
		Prefix: "StartDownRec",
		RemoteAddressHeaders: []string{"X-Forwarded-For"},
		OutputFlags: log.LstdFlags,
	})

	mux := http.NewServeMux()
	mux.Handle("/", loggerMiddleware.Handler(http.HandlerFunc(f.Exec)))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
