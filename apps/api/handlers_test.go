package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	api_utils "workout-tracker/libs/api/utils"
)

var app = &api_utils.Application{
	ErrorLog: log.New(io.Discard, "", 0),
	InfoLog:  log.New(io.Discard, "", 0),
}

const methodNotAllowedBody = "Method not allowed"
const wrongFormField = "The form does not contain any file under activity form field"

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

		uploadHandler(app)(rr, r)

		resp := rr.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
		}
		defer resp.Body.Close()
	})

	for _, tt := range testsMethodNotAllowed {
		t.Run(fmt.Sprint("Upload handler -> ", tt.name), func(t *testing.T) {
			rr := httptest.NewRecorder()
			// Initialize a new dummy http.Request.
			r, err := http.NewRequest(tt.method, "/upload-activity", nil)
			if err != nil {
				t.Fatal(err)
			}
			uploadHandler(app)(rr, r)

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

		uploadHandler(app)(rr, r)

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
