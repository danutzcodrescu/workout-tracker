package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	api_http "workout-tracker/libs/api-http"
	api_utils "workout-tracker/libs/api-utils"
)

const PORT = 8080

func setupRoutes(app *api_utils.Application) *http.ServeMux {
	// generate golang get route /ping with response OK
	mux := http.NewServeMux()
	api_http.Handlers(app, mux)

	// mux := http.NewServeMux()
	// mux.HandleFunc()
	// mux.HandleFunc("/upload-activity", uploadHandler(app))
	return mux
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &api_utils.Application{ErrorLog: errorLog, InfoLog: infoLog}
	mux := setupRoutes(app)
	log.Printf("Starting server at port %d\n", 8080)
	if err := http.ListenAndServe(":"+fmt.Sprintf("%d", PORT), mux); err != nil {
		log.Fatal(err)
	}
}
