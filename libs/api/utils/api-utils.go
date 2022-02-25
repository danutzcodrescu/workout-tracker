package api_utils

import (
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func ServerError(w http.ResponseWriter, err error, msg string) func(*Application) {
	return func(app *Application) {
		trace := fmt.Sprint(msg, "\n", err.Error())
		app.ErrorLog.Println(trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func ClientError(w http.ResponseWriter, err error, errorText string) func(*Application) {
	return func(app *Application) {
		app.ErrorLog.Println(err)
		http.Error(w, errorText, http.StatusBadRequest)
	}

}
