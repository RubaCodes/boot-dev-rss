package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rubacodes/boot-dev-rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// load enviroment variables
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Missing Server Configuration port")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Missing Database Configuration port")
	}
	// initialize sql connection
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to Database")
	}
	apiCfg := apiConfig{
		DB: database.New(conn),
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
	//Users
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser)) //get user with auth
	v1Router.Post("/users", apiCfg.handlerCreateUser)                    //create user
	//Feeds
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)                           // create feed
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed)) // create feed
	//FeedsFollow
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))                      // get feed_follow
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))                  // create feed_follow
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow)) // delete feed_follow

	// mount routes paths..
	router.Mount("/v1", v1Router)
	//create new server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Println("Server starting @ port:", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
