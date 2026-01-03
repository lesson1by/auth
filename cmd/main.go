package main

import (
	"authProject/internal/config"
	"authProject/internal/handlers"
	"authProject/internal/service"
	"authProject/internal/store"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	if err := store.CreateUsersTable(db); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	userStore := store.NewPgxStore(db)

	serv := service.NewUserService(cfg, userStore)

	loginHandler := &handlers.LoginHandlers{Serv: serv}
	http.HandleFunc("/login", loginHandler.Handle)

	registerHandler := &handlers.RegisterHandlers{Serv: serv}
	http.HandleFunc("/register", registerHandler.Handle)

	verifyHandler := &handlers.VerifyHandlers{Serv: serv}
	http.HandleFunc("/verify", verifyHandler.Handle)

	port := strconv.Itoa(cfg.Server.Port)

	log.Println("Server is running on port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
