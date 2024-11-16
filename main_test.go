package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up environment variables for testing
	os.Setenv("WEATHERAPI_KEY", "YOUR_WEATHERAPI_KEY")
	os.Exit(m.Run())
}

func TestValidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep/01310930", nil)
	rr := httptest.NewRecorder()

	handleCEPRequest(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}
	if tempC, ok := response["temp_C"]; !ok || tempC == nil {
		t.Errorf("handler returned JSON without temp_C field or value")
	}
	if tempF, ok := response["temp_F"]; !ok || tempF == nil {
		t.Errorf("handler returned JSON without temp_F field or value")
	}
	if tempK, ok := response["temp_K"]; !ok || tempK == nil {
		t.Errorf("handler returned JSON without temp_K field or value")
	}
}

func TestInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/cep/12345678", nil)
	rr := httptest.NewRecorder()

	handleCEPRequest(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	expectedBody := "Method not allowed\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestInvalidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep/invalidcep", nil)
	rr := httptest.NewRecorder()

	handleCEPRequest(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	expectedBody := "invalid zipcode\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestCEPNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep/00000000", nil)
	rr := httptest.NewRecorder()

	handleCEPRequest(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expectedBody := "can not find zipcode\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}
