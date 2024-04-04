package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes_UnmatchedRoute(t *testing.T) {
	router := NewRouter()

	// Insert your routes and handlers
	router.Insert("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		fmt.Fprintln(w, "home")
	}, "GET")

	router.Insert("/navdata", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		fmt.Fprintln(w, "navdata")
	}, "GET")

	router.Insert("/products/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		fmt.Fprintln(w, "product")
	}, "GET")

	// Test cases
	testCases := []struct {
		route    string
		expected int
	}{
		{"/", 200},
		{"/navdata", 200},
		{"/products/1", 200},
		{"/data", 404}, // Non-existent route
	}


	for _, testcase := range testCases {
		req, _ := http.NewRequest("GET", testcase.route, nil)
		rec := httptest.NewRecorder()
				
		router.Router(rec, req)

		if rec.Result().StatusCode != testcase.expected {
			t.Errorf("Expected '%d' but got '%d'", testcase.expected, rec.Result().StatusCode)
		}
	}


}

