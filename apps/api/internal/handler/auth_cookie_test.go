package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAuthLoginSetsCookie verifies that login sets the httpOnly cookie.
func TestAuthLoginSetsCookie(t *testing.T) {
	env := newAuthTestEnv(t)

	body := `{"username":"` + testSeedUsername + `","password":"` + testSeedPassword + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	env.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	cookies := rec.Result().Cookies()
	var found *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			found = c
			break
		}
	}
	if found == nil {
		t.Fatal("refresh_token cookie not set")
	}
	if !found.HttpOnly {
		t.Error("cookie should be httpOnly")
	}
	if found.Path != "/api/auth" {
		t.Errorf("expected Path=/api/auth, got %s", found.Path)
	}
	if found.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", found.SameSite)
	}
	if found.Value == "" {
		t.Error("cookie value is empty")
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	refreshToken, _ := resp["refreshToken"].(string)
	if refreshToken == "" {
		t.Error("JSON response should still contain refreshToken for non-browser clients")
	}
	user, _ := resp["user"].(map[string]any)
	if user["id"] != "user-admin" {
		t.Errorf("expected user user-admin, got %v", user["id"])
	}
}

// TestAuthRefreshThreeLayerFallback verifies Cookie → Header → Body priority.
func TestAuthRefreshThreeLayerFallback(t *testing.T) {
	env := newAuthTestEnv(t)

	tests := []struct {
		name         string
		useCookie    bool
		useHeader    bool
		useBody      bool
		expectStatus int
		expectError  string
	}{
		{
			name:         "cookie wins (priority 1)",
			useCookie:    true,
			useHeader:    false,
			useBody:      false,
			expectStatus: http.StatusOK,
		},
		{
			name:         "header wins when cookie empty (priority 2)",
			useCookie:    false,
			useHeader:    true,
			useBody:      false,
			expectStatus: http.StatusOK,
		},
		{
			name:         "body wins when cookie and header empty (priority 3)",
			useCookie:    false,
			useHeader:    false,
			useBody:      true,
			expectStatus: http.StatusOK,
		},
		{
			name:         "all empty returns 400",
			useCookie:    false,
			useHeader:    false,
			useBody:      false,
			expectStatus: http.StatusBadRequest,
			expectError:  "MISSING_REFRESH_TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Login fresh for each subtest to get a valid refresh token
			loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login",
				strings.NewReader(`{"username":"`+testSeedUsername+`","password":"`+testSeedPassword+`"}`))
			loginReq.Header.Set("Content-Type", "application/json")
			loginRec := httptest.NewRecorder()
			env.mux.ServeHTTP(loginRec, loginReq)

			if loginRec.Code != http.StatusOK {
				t.Fatalf("login failed: %d %s", loginRec.Code, loginRec.Body.String())
			}

			var loginResp map[string]any
			if err := json.NewDecoder(loginRec.Body).Decode(&loginResp); err != nil {
				t.Fatal(err)
			}
			refresh, _ := loginResp["refreshToken"].(string)
			if refresh == "" {
				t.Fatal("login did not return refreshToken")
			}

			// Build refresh request with token in the specified layer
			var body string
			if tt.useBody {
				body = `{"refreshToken":"` + refresh + `"}`
			} else {
				body = `{}`
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.useHeader {
				req.Header.Set("X-Refresh-Token", refresh)
			}
			if tt.useCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refresh})
			}
			rec := httptest.NewRecorder()

			env.mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectStatus, rec.Code, rec.Body.String())
			}
			if tt.expectError != "" && !strings.Contains(rec.Body.String(), tt.expectError) {
				t.Errorf("expected error code %s in body, got: %s", tt.expectError, rec.Body.String())
			}
			if tt.expectStatus == http.StatusOK {
				var resp map[string]any
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatal(err)
				}
				accessToken, _ := resp["accessToken"].(string)
				refreshToken, _ := resp["refreshToken"].(string)
				if accessToken == "" || refreshToken == "" {
					t.Error("expected new token pair in JSON response")
				}
				cookies := rec.Result().Cookies()
				var found *http.Cookie
				for _, c := range cookies {
					if c.Name == "refresh_token" {
						found = c
						break
					}
				}
				if found == nil {
					t.Error("refresh should update the cookie")
				} else if found.Value == "" {
					t.Error("cookie value should be the new refresh token")
				}
			}
		})
	}
}

// TestAuthLogoutClearsCookie verifies that logout clears the httpOnly cookie.
func TestAuthLogoutClearsCookie(t *testing.T) {
	env := newAuthTestEnv(t)

	// Login first
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"`+testSeedUsername+`","password":"`+testSeedPassword+`"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	env.mux.ServeHTTP(loginRec, loginReq)

	var loginResp map[string]any
	json.NewDecoder(loginRec.Body).Decode(&loginResp)
	refresh, _ := loginResp["refreshToken"].(string)

	// Logout with refresh token in body
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout",
		strings.NewReader(`{"refreshToken":"`+refresh+`"}`))
	logoutReq.Header.Set("Content-Type", "application/json")
	logoutRec := httptest.NewRecorder()

	env.mux.ServeHTTP(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", logoutRec.Code, logoutRec.Body.String())
	}

	cookies := logoutRec.Result().Cookies()
	var found *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			found = c
			break
		}
	}
	if found == nil {
		t.Fatal("refresh_token cookie should be present (with MaxAge=-1 to clear)")
	}
	if found.MaxAge != -1 {
		t.Errorf("expected MaxAge=-1 to clear cookie, got %d", found.MaxAge)
	}
}

// TestAuthLogoutViaCookie verifies logout works when token comes from cookie.
func TestAuthLogoutViaCookie(t *testing.T) {
	env := newAuthTestEnv(t)

	// Login first
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"`+testSeedUsername+`","password":"`+testSeedPassword+`"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	env.mux.ServeHTTP(loginRec, loginReq)

	var loginResp map[string]any
	json.NewDecoder(loginRec.Body).Decode(&loginResp)
	refresh, _ := loginResp["refreshToken"].(string)

	// Logout with refresh token in cookie only
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", strings.NewReader(`{}`))
	logoutReq.Header.Set("Content-Type", "application/json")
	logoutReq.AddCookie(&http.Cookie{Name: "refresh_token", Value: refresh})
	logoutRec := httptest.NewRecorder()

	env.mux.ServeHTTP(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", logoutRec.Code, logoutRec.Body.String())
	}
}
