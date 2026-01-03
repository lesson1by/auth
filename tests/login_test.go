package tests

import (

    "net/http"
    "net/http/httptest"
    "testing"

    "authProject/internal/handlers"
    "github.com/gojuno/minimock/v3"

)

func TestSuccessLogin(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    userServiceMock.ValidateCredentialsMock.Expect("username", "password").Return(true, nil)
    userServiceMock.GenerateTokenMock.Expect("username").Return("token", nil)

    h := &handlers.LoginHandlers{Serv: userServiceMock}
    req := httptest.NewRequest("POST", "/login", nil)
    req.SetBasicAuth("username", "password")
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status OK, got %v", resp.StatusCode)
    }
    if got := resp.Header.Get("Authorization"); got != "Bearer token" {
        t.Errorf("expected Authorization header 'Bearer token', got '%v'", got)
    }
}

func TestWrongPasswordLogin(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    userServiceMock.ValidateCredentialsMock.Expect("username", "wrongpass").Return(false, nil)
    h := &handlers.LoginHandlers{Serv: userServiceMock}
    req := httptest.NewRequest("POST", "/login", nil)
    req.SetBasicAuth("username", "wrongpass")
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()

    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("expected status Unauthorized, got %v", resp.StatusCode)
    }
}

func TestMissingOrWrongHeaderLogin(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    h := &handlers.LoginHandlers{Serv: userServiceMock}

    // Без Basic Auth
    req := httptest.NewRequest("POST", "/login", nil)
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()
    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("expected status Unauthorized (401) when no Basic Auth, got %v", resp.StatusCode)
    }

    // Bad method (например, GET вместо POST)
    req = httptest.NewRequest("GET", "/login", nil)
    rec = httptest.NewRecorder()
    h.Handle(rec, req)
    resp = rec.Result()
    if resp.StatusCode != http.StatusMethodNotAllowed {
        t.Errorf("expected status Method Not Allowed (405) when not POST, got %v", resp.StatusCode)
    }
}