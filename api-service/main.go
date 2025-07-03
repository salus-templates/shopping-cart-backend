package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time" // Added for http.Client timeout
)

// LoginRequest represents the structure of the incoming JSON request for login
type LoginRequest struct {
	Passkey string `json:"passkey"`
}

// LoginResponse represents the structure of the JSON response for login
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Product struct to match the structure of products from the Dotnet service
type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"` // Use float64 for decimal in Go
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
}

// authHandler handles authentication requests
func authHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow requests from any origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the configured passkey from an environment variable
	configuredPasskey := os.Getenv("AUTH_PASSKEY")
	if configuredPasskey == "" {
		log.Println("AUTH_PASSKEY environment variable is not set. Using default '12345'.")
		configuredPasskey = "12345" // Fallback for development if not set
	}

	// Compare the provided passkey with the configured passkey
	var resp LoginResponse
	if req.Passkey == configuredPasskey {
		resp = LoginResponse{Success: true, Message: "Authentication successful"}
		log.Printf("Login attempt for passkey '%s': SUCCESS", req.Passkey)
	} else {
		resp = LoginResponse{Success: false, Message: "Invalid passkey"}
		log.Printf("Login attempt for passkey '%s': FAILED (Incorrect passkey)", req.Passkey)
	}

	// Set content type and encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// productsHandler fetches, decodes, re-encodes, and responds with products
func productsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for any origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dotnetProductsApiURL := os.Getenv("DOTNET_PRODUCTS_API_URL")
	if dotnetProductsApiURL == "" {
		log.Println("DOTNET_PRODUCTS_API_URL environment variable is not set. Using default 'http://localhost:8080'.")
		dotnetProductsApiURL = "http://localhost:8080" // Default for development
	}

	// Construct the full URL for the Dotnet service
	targetURL := fmt.Sprintf("%s/all-products", dotnetProductsApiURL)
	log.Printf("Fetching products from Dotnet Products Service: %s", targetURL)

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(targetURL)
	if err != nil {
		log.Printf("Error fetching products from Dotnet service: %v", err)
		http.Error(w, "Failed to fetch products from backend service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Dotnet service returned non-OK status: %d", resp.StatusCode)
		http.Error(w, fmt.Sprintf("Backend service error: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	// Decode the JSON response from the Dotnet service
	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		log.Printf("Error decoding products from Dotnet service: %v", err)
		http.Error(w, "Failed to parse products data from backend", http.StatusInternalServerError)
		return
	}

	// --- This is where you can add logic to modify the 'products' slice if needed ---
	// Example: Add a new field, filter, sort, etc.
	// For now, we just re-encode it as is.
	// --------------------------------------------------------------------------------

	// Re-encode the products slice as JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Error encoding products for response: %v", err)
		// Note: Cannot send HTTP error after headers are written, just log.
	}
}

func main() {
	// Register the handlers
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/products", productsHandler) // Endpoint for products

	// Define the port to listen on
	port := "8090" // Default port for the Go app
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Go authentication and products processing proxy service listening on :%s\n", port)
	log.Printf("Go authentication and products processing proxy service starting on port %s", port)
	// Start the HTTP server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
