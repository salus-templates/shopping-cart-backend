package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestAuthHandler_Success tests successful authentication
func TestAuthHandler_Success(t *testing.T) {
	// Set a test passkey environment variable
	os.Setenv("AUTH_PASSKEY", "testpasskey")
	defer os.Unsetenv("AUTH_PASSKEY") // Clean up after test

	// Create a new HTTP request with a valid passkey
	loginReq := LoginRequest{Passkey: "testpasskey"}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the authHandler function
	authHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body
	var response LoginResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err == nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	// Check the response body
	if !response.Success {
		t.Errorf("handler returned unexpected success status: got %v want %v",
			response.Success, true)
	}
	expectedMessage := "Authentication successful"
	if response.Message != expectedMessage {
		t.Errorf("handler returned unexpected message: got %v want %v",
			response.Message, expectedMessage)
	}
}

// TestAuthHandler_Failure tests failed authentication due to incorrect passkey
func TestAuthHandler_Failure(t *testing.T) {
	// Set a test passkey environment variable
	os.Setenv("AUTH_PASSKEY", "testpasskey")
	defer os.Unsetenv("AUTH_PASSKEY") // Clean up after test

	// Create a new HTTP request with an invalid passkey
	loginReq := LoginRequest{Passkey: "wrongpasskey"}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the authHandler function
	authHandler(rr, req)

	// Check the status code (should still be 200 OK, but success: false in body)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body
	var response LoginResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err == nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	// Check the response body
	if response.Success {
		t.Errorf("handler returned unexpected success status: got %v want %v",
			response.Success, false)
	}
	expectedMessage := "Invalid passkey"
	if response.Message != expectedMessage {
		t.Errorf("handler returned unexpected message: got %v want %v",
			response.Message, expectedMessage)
	}
}

// TestAuthHandler_MethodNotAllowed tests handling of non-POST requests
func TestAuthHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth", nil) // Use GET method
	rr := httptest.NewRecorder()

	authHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for GET: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

// TestAuthHandler_InvalidJSON tests handling of invalid JSON body
func TestAuthHandler_InvalidJSON(t *testing.T) {
	reqBody := []byte(`{"passkey": "testpasskey", "extra": }`) // Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	authHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestAuthHandler_OptionsMethod tests handling of OPTIONS preflight request
func TestAuthHandler_OptionsMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/auth", nil)
	rr := httptest.NewRecorder()

	authHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for OPTIONS: got %v want %v",
			status, http.StatusOK)
	}
}
