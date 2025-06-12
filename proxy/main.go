package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/yaml.v2"
)

type RouteConfig struct {
	Routes map[string]string `yaml:"routes"`
}

var routeConfig RouteConfig

func loadRoutes() {
	log.Println("Loading routes configuration...")
	data, err := os.ReadFile("routes.yaml")
	if err != nil {
		log.Fatalf("Failed to read routing config: %v", err)
	}
	if err := yaml.Unmarshal(data, &routeConfig); err != nil {
		log.Fatalf("Failed to parse routing config: %v", err)
	}
	log.Printf("Loaded routes: %+v\n", routeConfig.Routes)
}

func reverseProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Printf("Error parsing target URL %s: %v", target, err)
		return nil
	}

	log.Printf("Creating reverse proxy for target: %s", target)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Add custom director for more control and logging
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		log.Printf("Proxying request: %s %s to %s", req.Method, req.URL.Path, target)
		originalDirector(req)
		log.Printf("Final request URL: %s", req.URL.String())
	}

	// Add error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for %s %s: %v", r.Method, r.URL.Path, err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	// Add response modifier for logging
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Received response: %d %s for %s", resp.StatusCode, resp.Status, resp.Request.URL.Path)
		return nil
	}

	return proxy
}

func validateJWT(r *http.Request) (jwt.MapClaims, error) {
	log.Printf("Validating JWT for request: %s %s", r.Method, r.URL.Path)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Missing Authorization header")
		return nil, fmt.Errorf("Missing Authorization header")
	}

	log.Printf("Authorization header present: %s", authHeader[:min(len(authHeader), 20)]+"...")

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		log.Println("Authorization header doesn't start with 'Bearer '")
		return nil, fmt.Errorf("Invalid authorization format")
	}

	log.Println("Reading public key...")
	pubKeyData, err := os.ReadFile("public.pem")
	if err != nil {
		log.Printf("Failed to read public key: %v", err)
		return nil, fmt.Errorf("Failed to read public key: %v", err)
	}

	log.Println("Parsing RSA public key...")
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		log.Printf("Invalid public key: %v", err)
		return nil, fmt.Errorf("Invalid public key: %v", err)
	}

	log.Println("Parsing JWT token...")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		log.Printf("JWT parsing error: %v", err)
		return nil, fmt.Errorf("Invalid JWT: %v", err)
	}

	if !token.Valid {
		log.Println("JWT token is invalid")
		return nil, fmt.Errorf("Invalid JWT token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid JWT claims format")
		return nil, fmt.Errorf("Invalid claims")
	}

	log.Printf("JWT validation successful. Claims: %+v", claims)
	return claims, nil
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("=== Incoming Request ===")
		log.Printf("Method: %s", r.Method)
		log.Printf("URL: %s", r.URL.String())
		log.Printf("Remote Addr: %s", r.RemoteAddr)
		log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))
		log.Printf("Content-Length: %d", r.ContentLength)

		// Log all headers (be careful with sensitive data)
		log.Println("Headers:")
		for key, values := range r.Header {
			for _, value := range values {
				if strings.ToLower(key) == "authorization" {
					log.Printf("  %s: %s", key, value[:min(len(value), 20)]+"...")
				} else {
					log.Printf("  %s: %s", key, value)
				}
			}
		}

		next(w, r)

		duration := time.Since(start)
		log.Printf("Request completed in %v", duration)
		log.Printf("=== End Request ===\n")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler called for: %s %s", r.Method, r.URL.Path)

	claims, err := validateJWT(r)
	if err != nil {
		log.Printf("JWT validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tier, ok := claims["tier"].(string)
	if !ok {
		log.Printf("Missing or invalid tier claim in JWT: %+v", claims)
		http.Error(w, "Missing tier claim", http.StatusBadRequest)
		return
	}

	log.Printf("Extracted tier from JWT: %s", tier)

	target, ok := routeConfig.Routes[tier]
	if !ok {
		log.Printf("Unknown tier '%s'. Available tiers: %+v", tier, routeConfig.Routes)
		http.Error(w, "Unknown tier", http.StatusForbidden)
		return
	}

	log.Printf("Routing request for tier '%s' to target '%s'", tier, target)

	proxy := reverseProxy(target)
	if proxy == nil {
		log.Printf("Failed to create reverse proxy for target: %s", target)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Executing reverse proxy...")
	proxy.ServeHTTP(w, r)
	log.Printf("Reverse proxy execution completed")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check requested from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting reverse proxy server...")

	loadRoutes()

	// Add health check endpoint
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", loggingMiddleware(handler))

	log.Println("Reverse proxy listening on :9000")
	log.Println("Health check available at /health")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
