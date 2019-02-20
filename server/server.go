package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/lordpuma/bettershifts-backend"
	"github.com/lordpuma/bettershifts-backend/auth"
  "github.com/go-chi/chi"
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

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(bettershifts.NewExecutableSchema(bettershifts.Config{Resolvers: &resolver})))

	err = http.ListenAndServe(port, router)
	if err != nil {
		panic(err)
	}



	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
