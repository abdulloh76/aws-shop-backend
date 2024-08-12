package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	r := mux.NewRouter()

	// Route to handle all requests
	r.PathPrefix("/").HandlerFunc(redirectHandler)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	targetService := getTargetService(r)

	if targetService == "" {
		http.Error(w, "Cannot process request", http.StatusBadGateway)
		return
	}

	redirectURL := targetService + r.RequestURI

	log.Printf("Redirecting to: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func getTargetService(r *http.Request) string {
	// Example logic to choose the target service based on request URI or headers
	if r.RequestURI == "/product" || r.Header.Get("X-Service") == "ProductService" {
		return os.Getenv("PRODUCT_SERVICE_API_URL")
	} else if r.RequestURI == "/import" || r.Header.Get("X-Service") == "ImportService" {
		return os.Getenv("IMPORT_SERVICE_API_URL")
	} else if r.RequestURI == "/cart" || r.Header.Get("X-Service") == "CartService" {
		return os.Getenv("CART_SERVICE_API_URL")
	}

	return ""
}
