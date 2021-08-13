package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)


func TestMain(m *testing.M) {
	flag.Parse()
	http.DefaultServeMux.HandleFunc("/files", getSummary)
	http.DefaultServeMux.HandleFunc("/statistics", getSummary)

	os.Exit(m.Run())
}

func TestHandleGetSummary(t *testing.T) {
	t.Run("get files", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/files?path=var", nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Response code is %v", result.StatusCode)
		}

		res := new(FilesResponse)
		_ = json.Unmarshal(w.Body.Bytes(), res)
		expected := "C:\\var"
		if res.Path != expected {
			t.Errorf("expected: %s, actual: %s", expected, res.Path)
		}
		t.Logf("files: %+v", res)
	})

	t.Run("get statistics", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/statistics?path=var", nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)

		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Response code is %v", result.StatusCode)
		}

		res := new(StatisticsResponse)
		_ = json.Unmarshal(w.Body.Bytes(), res)
		expected := "C:\\var"
		if res.Path != expected {
			t.Errorf("expected: %s, actual: %s", expected, res.Path)
		}
		t.Logf("statistics: %+v", res)
	})
}
