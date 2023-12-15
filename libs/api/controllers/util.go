package api_controllers

import (
	"fmt"
	"log"
	"net/http"
	api_repositories "workout-tracker/libs/api/repositories"
)

type Application struct {
	ErrorLog     *log.Logger
	InfoLog      *log.Logger
	Repositories api_repositories.Repositories
}

func serverError(w http.ResponseWriter, err error, msg string) func(*Application) {
	return func(app *Application) {
		trace := fmt.Sprint(msg, "\n", err.Error())
		app.ErrorLog.Println(trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func clientError(w http.ResponseWriter, err error, errorText string) func(*Application) {
	return func(app *Application) {
		app.ErrorLog.Println(err)
		http.Error(w, errorText, http.StatusBadRequest)
	}

}
