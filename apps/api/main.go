package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	api_utils "workout-tracker/libs/api/utils"
)

const port = 8080

func setupRoutes(app *api_utils.Application) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload-activity", uploadHandler(app))
	return mux
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &api_utils.Application{ErrorLog: errorLog, InfoLog: infoLog}
	mux := setupRoutes(app)
	log.Printf("Starting server at port %d\n", port)
	if err := http.ListenAndServe(":"+fmt.Sprintf("%d", port), mux); err != nil {
		log.Fatal(err)
	}
}
