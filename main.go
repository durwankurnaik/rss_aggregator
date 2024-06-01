package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // used to bring the env vars from .env file to here

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not bound in the environment")
	}

	router := chi.NewRouter()

	// cors config to allow the requests in our browser to connect with the backend
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)

	router.Mount("/v1", v1Router) // good practice to have this setup, so that you can migrate to v2 easily

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe() // this line creates the process which runs the go server
	if err != nil {
		log.Fatal(err) // in idea scenario, this won't hit, that means good thing for the dev
	}
}
