package api_utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app = &Application{
	ErrorLog: log.New(io.Discard, "", 0),
	InfoLog:  log.New(io.Discard, "", 0),
}

func TestServerError(t *testing.T) {

	tests := []struct {
		name     string
		wantCode int
		wantBody []byte
	}{
		{"Generic server error", http.StatusInternalServerError, []byte(fmt.Sprint((http.StatusText(http.StatusInternalServerError)), "\n"))},
		{"Other error", http.StatusInternalServerError, []byte(fmt.Sprint(http.StatusText(http.StatusInternalServerError), "\n"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ServerError(w, errors.New("test error"), "test error body")(app)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, res.StatusCode)
			}
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != string(tt.wantBody) {
				t.Errorf("want body to equal %q; got %q", tt.wantBody, string(data))
			}
		})
	}
}

func TestClientError(t *testing.T) {
	tests := []struct {
		name     string
		wantCode int
		wantBody []byte
	}{
		{"Upload file", http.StatusBadRequest, []byte("The uploaded file is too big. Please choose an file that's less than")},
		{"Form error", http.StatusBadRequest, []byte("The form does not contain any file under activity form field")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ClientError(w, errors.New("test error"), string(tt.wantBody))(app)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, res.StatusCode)
			}
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != fmt.Sprint(string((tt.wantBody)), "\n") {
				t.Errorf("want body to equal %q; got %q", tt.wantBody, string(data))
			}
		})
	}
}
