package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tests = []struct {
	method  string
	route   string
	body    io.Reader
	handler handlerT
	want    string
}{
	{"GET", "/data", nil, dataHandler, "hi"},
}

func TestRoutes(t *testing.T) {
	ts := httptest.NewServer(newMux())
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.route, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+tc.route, tc.body)
			if err != nil {
				t.Fatal(err.Error())
			}

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err.Error())
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err.Error())
			}

			if string(body) != tc.want {
				t.Fatalf("got %q; want %q", body, tc.want)
			}
		})
	}
}
