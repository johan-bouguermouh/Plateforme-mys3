package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPermissionMiddleware vérifie les différentes conditions d'accès
func TestPermissionMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name               string
		path               string
		user               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "Public read access",
			path:               "/bucket1/file",
			user:               "user1",
			authHeader:         "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Private access with valid auth",
			path:               "/bucket1/file",
			user:               "user2",
			authHeader:         "Bearer mysecret",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Private access with invalid auth",
			path:               "/bucket1/file",
			user:               "user2",
			authHeader:         "Bearer wrongsecret",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "User not in ACL",
			path:               "/bucket1/file",
			user:               "user5",
			authHeader:         "",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "Invalid path",
			path:               "/invalidpath/file",
			user:               "user1",
			authHeader:         "",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("X-User", tt.user)
			req.Header.Set("Authorization", tt.authHeader)

			rr := httptest.NewRecorder()
			middleware := PermissionMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
		})
	}
}
