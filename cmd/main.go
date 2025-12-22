package main

import (
	"authProject/internal/config"
	"authProject/internal/handlers"
	"authProject/internal/service"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
		
	serv := service.NewUserService(cfg)

	loginHandler := &handlers.LoginHandlers{Serv: serv}
	http.HandleFunc("/login", loginHandler.Handle)

	registerHandler := &handlers.RegisterHandlers{Serv: serv}
	http.HandleFunc("/register", registerHandler.Handle)

	verifyHandler := &handlers.VerifyHandlers{Serv: serv}
	http.HandleFunc("/verify", verifyHandler.Handle)

	port := strconv.Itoa(cfg.Server.Port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running on port", port)
}
