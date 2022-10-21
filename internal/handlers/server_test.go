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
		code     int
		response string
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
				code:     400,
				response: "Only text/plain Content-Type header is allowed",
			},
		},
		{
			name:        "positive",
			method:      http.MethodPost,
			args:        "gauge/gAlloc/1.23445",
			contentType: "text/plain",
			want: want{
				code:     200,
				response: "data accepted",
			},
		},
		{
			name:        "wrong args",
			method:      http.MethodPost,
			args:        "gauge/gAlloc/",
			contentType: "text/plain",
			want: want{
				code:     400,
				response: "Incorrect metric value",
			},
		},
		{
			name:        "no metric name (gauge)",
			method:      http.MethodPost,
			args:        "gauge/",
			contentType: "text/plain",
			want: want{
				code:     404,
				response: "metric name is not found",
			},
		},
		{
			name:        "no metric name #2 (counter)",
			method:      http.MethodPost,
			args:        "counter/",
			contentType: "text/plain",
			want: want{
				code:     404,
				response: "metric name is not found",
			},
		},
		{
			name:        "wrong metric value #1 (gauge)",
			method:      http.MethodPost,
			args:        "gauge/gAlloc/none",
			contentType: "text/plain",
			want: want{
				code:     400,
				response: "Incorrect metric value",
			},
		},
		{
			name:        "wrong metric value #2 (counter)",
			method:      http.MethodPost,
			args:        "counter/cAlloc/none",
			contentType: "text/plain",
			want: want{
				code:     400,
				response: "Incorrect metric value",
			},
		},
		{
			name:        "counter update",
			method:      http.MethodPost,
			args:        "counter/Alloc/100",
			contentType: "text/plain",
			want: want{
				code:     200,
				response: "data accepted",
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

			if res.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
				t.Errorf("Expected Content-Type %s, got %s", "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			}
		})
	}
}
