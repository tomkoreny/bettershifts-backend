package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"os"
	"github.com/99designs/gqlgen/handler"
	"github.com/lordpuma/bettershifts"
	"github.com/lordpuma/bettershifts/auth"
  "github.com/go-chi/chi" 
	"github.com/rs/cors"
)

const defaultPort = "8080"


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
  /*	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Todo{})
	db.AutoMigrate(&models.Shift{})
	db.AutoMigrate(&models.Benefit{})
	db.AutoMigrate(&models.Workplace{}) **/

	var resolver = bettershifts.Resolver{Db: db}

  router := chi.NewRouter()

	router.Use(auth.Middleware(db))
  router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(bettershifts.NewExecutableSchema(bettershifts.Config{Resolvers: &resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
