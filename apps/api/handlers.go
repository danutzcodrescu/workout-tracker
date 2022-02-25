package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api_activity "workout-tracker/libs/api/activity"
	api_utils "workout-tracker/libs/api/utils"
)

const file_size_in_mb = 3

// 3MB file upload size
const max_upload_size = 1024 * 1024 * file_size_in_mb

func uploadHandler(app *api_utils.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, max_upload_size)
		if err := r.ParseMultipartForm(max_upload_size); err != nil {
			api_utils.ClientError(w, err, fmt.Sprintln("The uploaded file is too big. Please choose an file that's less than", file_size_in_mb, "MB in size."))(app)
			return
		}

		file, _, err := r.FormFile("activity")
		if err != nil {
			api_utils.ClientError(w, err, fmt.Sprintln("The form does not contain any file under activity form field"))(app)
			return
		}
		defer file.Close()
		fileBytes, _ := ioutil.ReadAll(file)

		activityWorkout, err := api_activity.ParseActivityFile(fileBytes)
		if err != nil {
			api_utils.ServerError(w, err, "Error parsing workout")(app)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activityWorkout)
	}
}
