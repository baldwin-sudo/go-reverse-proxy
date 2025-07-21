package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baldwin-sudo/go-reverse-proxy/internal/proxy"
)

func TestHandleReqOverProxy(t *testing.T) {
	// Step 1: Create a mock backend server
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("mock-backend1"))
	}))
	defer backend1.Close()

	// Step 2: Create a RouteConfig
	routeConfig := proxy.RouteConfig{
		Path:     "/test",
		Backends: []string{backend1.URL},
	}

	// Step 3: Create the Route from config
	route, err := proxy.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("Failed to initialize route: %v", err)
	}

	// Step 4: Create the handler
	handler := route.HandleReqOverProxy()

	// Step 5: Simulate a request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTeapot {
		t.Errorf("Expected status 418, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "mock-backend1" {
		t.Errorf("Expected body 'mock-backend1', got '%s'", string(body))
	}
}

func TestHandleReqOverProxyMultipleBackends(t *testing.T) {
	// Step 1: Create two mock backends
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("mock-backend1"))
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("mock-backend2"))
	}))
	defer backend2.Close()

	// Step 2: Create a RouteConfig with both backends
	routeConfig := proxy.RouteConfig{
		Path:     "/test",
		Backends: []string{backend1.URL, backend2.URL},
	}

	// Step 3: Initialize Route
	route, err := proxy.NewRoute(routeConfig)

	if err != nil {
		t.Fatalf("Failed to initialize route: %v", err)
	}

	handler := route.HandleReqOverProxy()

	// Step 4: Make requests and assert round-robin behavior
	expectedBodies := []string{"mock-backend1", "mock-backend2"}
	for _, expected := range expectedBodies {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if res.StatusCode != http.StatusTeapot {
			t.Errorf("Expected status 418, got %d", res.StatusCode)
		}

		if string(body) != expected {
			t.Errorf("Expected body '%s', got '%s'", expected, string(body))
		}
	}
}
