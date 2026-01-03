package tests


import (
	"errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "authProject/internal/handlers"
    "github.com/gojuno/minimock/v3"

)

func TestValidToken(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    userServiceMock.RefreshTokenMock.Expect("Bearer validToken").Return("newValidToken", nil) // валидный токен

    h := &handlers.VerifyHandlers{Serv: userServiceMock}
    req := httptest.NewRequest("POST", "/verify", nil)
    req.Header.Set("Authorization", "Bearer validToken")
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status OK, got %v", resp.StatusCode)
    }
    if got := resp.Header.Get("Authorization"); got != "newValidToken" {
        t.Errorf("expected Authorization header 'newValidToken', got '%v'", got)
    }
}

func TestExpiredToken(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    userServiceMock.RefreshTokenMock.Expect("Bearer expiredToken").Return("", errors.New("token expired")) // просроченный токен

    h := &handlers.VerifyHandlers{Serv: userServiceMock}
    req := httptest.NewRequest("POST", "/verify", nil)
    req.Header.Set("Authorization", "Bearer expiredToken")
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()
    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("expected status Unauthorized (401), got %v", resp.StatusCode)
    }
}

func TestFakeToken(t *testing.T) {
    mc := minimock.NewController(t)
    defer mc.Finish()
    userServiceMock := NewUserServiceMock(mc)
    userServiceMock.RefreshTokenMock.Expect("Bearer fakeToken").Return("", errors.New("invalid")) // поддельный токен

    h := &handlers.VerifyHandlers{Serv: userServiceMock}
    req := httptest.NewRequest("POST", "/verify", nil)
    req.Header.Set("Authorization", "Bearer fakeToken")
    rec := httptest.NewRecorder()
    h.Handle(rec, req)
    resp := rec.Result()
    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("expected status Unauthorized (401), got %v", resp.StatusCode)
    }
}