package handlers

import (
	"authProject/internal/service"
	"net/http"
	"strings"
)

type VerifyHandlers struct {
	Serv service.UserService
}

func (v *VerifyHandlers) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "missing bearer token", http.StatusUnauthorized)
		return
	}

	newToken, err := v.Serv.RefreshToken(authHeader)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", newToken)
	w.WriteHeader(http.StatusOK)
}
