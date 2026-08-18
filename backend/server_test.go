package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluateEndpoint_OK(t *testing.T) {
	srv := NewServer(nil)
	body := `{"cartonPricePerUnit":35,"loosePricePerUnit":40,"cartonSize":24,"spoilagePct":2,"capitalCostRatePct":24,"weeklySellThrough":24}`
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var r Result
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if r.Winner != "carton" {
		t.Fatalf("winner=%s want carton", r.Winner)
	}
}

func TestEvaluateEndpoint_ValidationError(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(`{"cartonSize":0,"weeklySellThrough":6,"loosePricePerUnit":1}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d want 400", rec.Code)
	}
}

func TestHistoryEndpoint_UnavailableWithoutDB(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/history", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d want 503", rec.Code)
	}
}
