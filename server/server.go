package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"os" 
  "github.com/gobuffalo/packr/v2"
  "github.com/rubenv/sql-migrate"

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

  box := packr.New("migrations", "./migrations")
  migrations := &migrate.PackrMigrationSource{
    Box: box,
  }
  n, err := migrate.Exec(db.DB(), "postgres", migrations, migrate.Up)
if err != nil {
  panic(err)
    // Handle errors!
}
fmt.Printf("Applied %d migrations!\n", n)

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
