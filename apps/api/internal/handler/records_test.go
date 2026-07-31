package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func recordsMux() *http.ServeMux {
	mux := http.NewServeMux()
	recordsHandler(mux)
	return mux
}

func getJSON(t *testing.T, mux *http.ServeMux, path string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, body
}

func TestRecordsListDefault(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	if len(items) != 8 {
		t.Fatalf("len(items) = %d, want 8", len(items))
	}
	if body["total"] != float64(8) {
		t.Fatalf("total = %v, want 8", body["total"])
	}
	if body["pageSize"] != float64(10) {
		t.Fatalf("pageSize = %v, want 10", body["pageSize"])
	}
	// Default sort is name asc: first item is "Acme Console".
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("first item = %v, want object", items[0])
	}
	if first["name"] != "Acme Console" {
		t.Fatalf("first name = %v, want Acme Console", first["name"])
	}
}

func TestRecordsListSearch(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records?q=alice")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	for _, item := range items {
		rec := item.(map[string]any)
		if rec["owner"] != "alice" && !strings.Contains(rec["name"].(string), "alice") {
			t.Fatalf("item %v does not match q=alice", rec)
		}
	}
}

func TestRecordsListSortDesc(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records?sort=updatedAt&order=desc")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	last, _ := items[0].(map[string]any)
	if last["name"] != "Globex Admin" {
		t.Fatalf("first (desc) name = %v, want Globex Admin", last["name"])
	}
}

func TestRecordsListPagination(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records?page=2&pageSize=3")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("items = %v, want array", body["items"])
	}
	if len(items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(items))
	}
	if body["total"] != float64(8) {
		t.Fatalf("total = %v, want 8", body["total"])
	}
	first, _ := items[0].(map[string]any)
	if first["name"] != "Initech Reports" {
		t.Fatalf("page 2 first name = %v, want Initech Reports", first["name"])
	}
}

func TestRecordsListInvalidParamsFailClosed(t *testing.T) {
	for _, path := range []string{
		"/api/records?sort=unknown",
		"/api/records?order=up",
		"/api/records?page=0",
		"/api/records?pageSize=abc",
	} {
		code, body := getJSON(t, recordsMux(), path)
		if code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d, want %d", path, code, http.StatusBadRequest)
		}
		if _, ok := body["error"]; !ok {
			t.Fatalf("%s: error code missing in %v", path, body)
		}
	}
}

func TestRecordsDetail(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records/rec-3")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d", code, http.StatusOK)
	}
	if body["id"] != "rec-3" {
		t.Fatalf("id = %v, want rec-3", body["id"])
	}
}

func TestRecordsDetailNotFound(t *testing.T) {
	code, body := getJSON(t, recordsMux(), "/api/records/rec-999")
	if code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", code, http.StatusNotFound)
	}
	if body["error"] != "RECORD_NOT_FOUND" {
		t.Fatalf("error = %v, want RECORD_NOT_FOUND", body["error"])
	}
}
