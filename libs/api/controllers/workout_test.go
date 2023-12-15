package api_controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	api_utils "workout-tracker/libs/api/utils"
)

var application = &Application{
	ErrorLog: log.New(io.Discard, "", 0),
	InfoLog:  log.New(io.Discard, "", 0),
}

const methodNotAllowedBody = "Method not allowed"
const wrongFormField = "The form does not contain any file under activity form field"

var expectedWorkout = &api_utils.Workout{Date: "2022-02-24T13:39:00Z", Laps: []api_utils.WorkoutLap{
	{
		StartTime:        "2022-02-24T13:39:00Z",
		TotalTimeSeconds: 60,
		DistanceMeters:   244,
		Calories:         16,
		Intensity:        "Active",
		Efforts: []api_utils.Effort{
			{
				Time:           1,
				DistanceMeters: 4,
				Cadence:        0,
				Watts:          61,
			},
			{
				Time:           2,
				DistanceMeters: 8,
				Cadence:        0,
				Watts:          61,
			},
			{
				Time:           58,
				DistanceMeters: 238,
				Cadence:        22,
				Watts:          201,
			},
		},
	},
	{
		StartTime:        "2022-02-24T13:40:00Z ",
		TotalTimeSeconds: 60,
		DistanceMeters:   172,
		Calories:         9,
		Intensity:        "Resting",
		Efforts: []api_utils.Effort{
			{
				Time:           2,
				DistanceMeters: 11,
				Cadence:        22,
				Watts:          201,
			},
			{
				Time:           58,
				DistanceMeters: 172,
				Cadence:        22,
				Watts:          63,
			},
		},
	},
}}

func TestUploadActivity(t *testing.T) {
	workout, err := os.Open("../../tools/testFiles/workout.tcx")
	if err != nil {
		t.Errorf("cannot parse workout file")
	}
	defer workout.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("activity", "../../tools/testFiles/workout.tcx")
	io.Copy(part, workout)
	writer.Close()

	testsMethodNotAllowed := []struct {
		name     string
		method   string
		wantCode int
		wantBody []byte
	}{
		{"Method not allowed GET", http.MethodGet, http.StatusMethodNotAllowed, []byte(methodNotAllowedBody)},
		{"Method not allowed PUT", http.MethodGet, http.StatusMethodNotAllowed, []byte(methodNotAllowedBody)},
		{"Method not allowed PATCH", http.MethodGet, http.StatusMethodNotAllowed, []byte(methodNotAllowedBody)},
		{"Method not allowed DELETE", http.MethodGet, http.StatusMethodNotAllowed, []byte(methodNotAllowedBody)},
	}

	t.Run("Upload handler - smoke test", func(t *testing.T) {
		rr := httptest.NewRecorder()
		// Initialize a new dummy http.Request.
		r, err := http.NewRequest(http.MethodPost, "/upload-activity", body)
		r.Header.Add("Content-type", writer.FormDataContentType())
		if err != nil {
			t.Fatal(err)
		}

		UploadActivityController(application)(rr, r)

		resp := rr.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
		}
		defer resp.Body.Close()
		result := api_utils.Workout{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalln(err)
		}
		if reflect.DeepEqual(result, expectedWorkout) {
			t.Errorf("want %+v; got %+v", expectedWorkout, result)
		}
	})

	for _, tt := range testsMethodNotAllowed {
		t.Run(fmt.Sprint("Upload handler -> ", tt.name), func(t *testing.T) {
			rr := httptest.NewRecorder()
			// Initialize a new dummy http.Request.
			r, err := http.NewRequest(tt.method, "/upload-activity", nil)
			if err != nil {
				t.Fatal(err)
			}
			UploadActivityController(application)(rr, r)

			resp := rr.Result()

			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != fmt.Sprint(string(tt.wantBody), "\n") {
				t.Errorf("wanted %s; got %s", methodNotAllowedBody, string(body))
			}
			if resp.Header.Get("Allow") != http.MethodPost {
				t.Errorf("wanted %s; got %s", http.MethodPost, resp.Header.Get("Allow"))
			}
		})
	}

	t.Run("Upload handler - fail on wrong form field submitted", func(t *testing.T) {
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "../../tools/testFiles/workout.tcx")
		io.Copy(part, workout)
		writer.Close()
		rr := httptest.NewRecorder()
		// Initialize a new dummy http.Request.
		r, err := http.NewRequest(http.MethodPost, "/upload-activity", body)
		r.Header.Add("Content-type", writer.FormDataContentType())
		if err != nil {
			t.Fatal(err)
		}

		UploadActivityController(application)(rr, r)

		resp := rr.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("want %d; got %d", http.StatusBadRequest, resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), wrongFormField) {
			t.Errorf("wanted %s; got %s", wrongFormField, string(body))
		}
	})

}
