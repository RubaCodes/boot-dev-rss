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
	// load enviroment variables
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Missing Server Configuration port")
	}
	//create a new router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// create a sub path router
	v1Router := chi.NewRouter()
	v1Router.Get("/healtz", handlerReadiness) //health check
	v1Router.Get("/err", handlerErr)          //err check
	// mount routes paths..
	router.Mount("/v1", v1Router)
	//create new server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Println("Server starting @ port:", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
