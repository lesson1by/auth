package handlers

import (
	"authProject/internal/service"
	"net/http"
)

type LoginHandlers struct {
	Serv service.UserService
} // спросить почему нельзя эту структура положить в /models и почему из-за этого получается цикл импортов
func (h *LoginHandlers) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="login"`)
		http.Error(w, "missing or invalid basic auth", http.StatusUnauthorized)
		return
	}

	ok, err := h.Serv.ValidateCredentials(username, password)
	if err != nil || !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="login"`)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	jwtToken, err := h.Serv.GenerateToken(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwtToken)
	w.WriteHeader(http.StatusOK)
}
