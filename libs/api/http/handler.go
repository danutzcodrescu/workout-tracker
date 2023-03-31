package api_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	api_controllers "workout-tracker/libs/api/controllers"
	api_utils "workout-tracker/libs/api/utils"
)

const FILE_SIZE_IN_MB = 3

// 3MB file upload size
const MAX_UPLOAD_SIZE = 1024 * 1024 * FILE_SIZE_IN_MB

func uploadHandler(app *api_utils.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			api_utils.ClientError(w, err, fmt.Sprintln("The uploaded file is too big. Please choose an file that's less than", FILE_SIZE_IN_MB, "MB in size."))(app)
			return
		}

		file, _, err := r.FormFile("activity")
		if err != nil {
			api_utils.ClientError(w, err, fmt.Sprintln("The form does not contain any file under activity form field"))(app)
			return
		}
		defer file.Close()
		fileBytes, _ := io.ReadAll(file)
		activityWorkout, err := api_controllers.ParseActivityFile(fileBytes)
		if err != nil {
			api_utils.ServerError(w, err, "Error parsing workout")(app)
		}
		err = app.Repositories.Activity.Insert(activityWorkout)
		if err != nil {
			api_utils.ServerError(w, err, "Error inserting record")(app)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activityWorkout)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`OK`))
}

func Handlers(app *api_utils.Application, mux *http.ServeMux) {
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/upload-activity", uploadHandler(app))
}
