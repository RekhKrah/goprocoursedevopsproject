package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string
		method      string
		contentType string
		args        string
		want        want
	}{
		{
			name:        "incorrect content-type",
			method:      http.MethodPost,
			contentType: "application/json",
			args:        "",
			want: want{
				code:        400,
				response:    "Only text/plain Content-Type header is allowed",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "wrong method",
			method:      http.MethodGet,
			args:        "gauge/Alloc/1.23445",
			contentType: "text/plain",
			want: want{
				code:        405,
				response:    "Only POST request are allowed",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "wrong method",
			method:      http.MethodPost,
			args:        "gauge/Alloc/1.23445",
			contentType: "text/plain",
			want: want{
				code:        200,
				response:    "data accepted",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "wrong args",
			method:      http.MethodPost,
			args:        "gauge/Alloc/",
			contentType: "text/plain",
			want: want{
				code:        400,
				response:    "Can't parse url '/update/gauge/Alloc/'",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, fmt.Sprintf("/update/%v", tt.args), nil)
			request.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateMetrics)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(string(resBody)) != tt.want.response {
				t.Errorf("Expected body '%s', got '%s'", tt.want.response, w.Body.String())
			}

			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
