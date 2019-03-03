package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"os" 
  _ "github.com/lib/pq"
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/99designs/gqlgen/handler"
	"github.com/lordpuma/bettershifts-backend"
	"github.com/lordpuma/bettershifts-backend/auth"
  "github.com/go-chi/chi"
)

const defaultPort = "8080"
const defaultDb = ""


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	dbAddr := os.Getenv("DATABASE_URL")
	if dbAddr == "" {
		dbAddr = defaultDb
	}

	db, err := gorm.Open("postgres", dbAddr)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
  driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
  m, err := migrate.NewWithDatabaseInstance(
      "file:///migrations",
      "postgres", driver)
  m.Steps(1)

	var resolver = bettershifts.Resolver{Db: db}

  router := chi.NewRouter()
	router.Use(auth.Middleware(db))

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(bettershifts.NewExecutableSchema(bettershifts.Config{Resolvers: &resolver})))

  err = http.ListenAndServe(":" + port, router)
	if err != nil {
		panic(err)
	}


	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
