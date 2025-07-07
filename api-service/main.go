package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

// Product struct to match the structure of products from the Dotnet service (now includes Stock)
type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"` // New: Stock quantity
}

// OrderItemRequest from React app
type OrderItemRequest struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// PlaceOrderRequest from React app to Go
type PlaceOrderRequest struct {
	Items           []OrderItemRequest `json:"items"`
	TotalAmount     float64            `json:"totalAmount"`
	DeliveryAddress string             `json:"deliveryAddress"`
	OrderDate       string             `json:"orderDate"`
}

// PlaceOrderResponse from Dotnet to Go, and then Go to React
type PlaceOrderResponse struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message,omitempty"`
	OrderId         string   `json:"orderId,omitempty"`
	OutOfStockItems []string `json:"outOfStockItems,omitempty"` // New: List of items that caused failure
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
		log.Printf("Error: Dotnet service returned non-OK status: %d", resp.StatusCode)
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

	// Re-encode the products slice as JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Error encoding products for response: %v", err)
	}
}

// orderHandler proxies and processes order requests to the Dotnet products-service
func orderHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for any origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dotnetProductsApiURL := os.Getenv("DOTNET_PRODUCTS_API_URL")
	if dotnetProductsApiURL == "" {
		log.Println("DOTNET_PRODUCTS_API_URL environment variable is not set. Using default 'http://localhost:8080'.")
		dotnetProductsApiURL = "http://localhost:8080" // Default for development
	}

	// Construct the full URL for the Dotnet service's place-order endpoint
	targetURL := fmt.Sprintf("%s/place-order", dotnetProductsApiURL)
	log.Printf("Proxying order request to Dotnet Products Service: %s", targetURL)

	// Decode the incoming order request from React
	var orderRequest PlaceOrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		log.Printf("Error decoding order request from client: %v", err)
		http.Error(w, "Invalid order request body", http.StatusBadRequest)
		return
	}

	// Re-encode the order request to send to Dotnet service
	requestBodyBytes, err := json.Marshal(orderRequest)
	if err != nil {
		log.Printf("Error marshalling order request for Dotnet: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a new HTTP POST request to the Dotnet service
	client := &http.Client{Timeout: 10 * time.Second}
	proxyReq, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Printf("Error creating proxy order request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	proxyReq.Header.Set("Content-Type", "application/json") // Ensure JSON content type for Dotnet

	// Perform the request to Dotnet
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error placing order with Dotnet service: %v", err)
		http.Error(w, "Failed to place order with backend service", http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	if code := proxyResp.StatusCode; code != http.StatusOK {
		log.Printf("Error: Dotnet service returned non-OK status: %d", code)
		http.Error(w, fmt.Sprintf("Backend service error: %d", code), http.StatusBadGateway)
		return
	}

	// Decode the response from the Dotnet service
	var orderResponse PlaceOrderResponse
	err = json.NewDecoder(proxyResp.Body).Decode(&orderResponse)
	if err != nil {
		log.Printf("Error decoding order response from Dotnet service: %v", err)
		http.Error(w, "Failed to parse order response from backend", http.StatusInternalServerError)
		return
	}

	// --- This is where you can add logic to modify the 'orderResponse' if needed ---
	// For now, we just re-encode it as is.
	// --------------------------------------------------------------------------------

	// Re-encode the Dotnet response and send it back to React
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(proxyResp.StatusCode) // Pass through the status code from Dotnet
	if err := json.NewEncoder(w).Encode(orderResponse); err != nil {
		log.Printf("Error encoding order response for client: %v", err)
	}
}

func main() {
	// Register the handlers
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/order", orderHandler) // New endpoint for order processing

	// Define the port to listen on
	port := "8080" // Default port for the Go app
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Go authentication, products and order processing proxy service listening on :%s\n", port)
	log.Printf("Go authentication, products and order processing proxy service starting on port %s", port)
	// Start the HTTP server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
