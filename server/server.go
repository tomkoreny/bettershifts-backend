package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"os"
  "context"

	"github.com/99designs/gqlgen/handler"
	"github.com/lordpuma/bettershifts-backend"
	"github.com/lordpuma/bettershifts-backend/models"
  "github.com/go-chi/chi"
)

const defaultPort = "8080"

var userCtxKey = &contextKey{"user"}
type contextKey struct {
	name string
}

func validateAndGetUserID(c *http.Cookie, db *gorm.DB) (string, error) {
  return "1", nil; 
}

func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("auth-cookie")

			// Allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			userId, err := validateAndGetUserID(c, db)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			// get the user from the database
			user := getUserByID(db, userId)

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(userCtxKey).(*models.User)
	return raw
}

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

	router.Use(Middleware(db))
  router := chi.NewRouter()

	router.Use(auth.Middleware(db))

	router.Handle("/", handler.Playground("Starwars", "/query"))
	router.Handle("/query",
		handler.GraphQL(starwars.NewExecutableSchema(starwars.NewResolver())),
	)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}


	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(bettershifts.NewExecutableSchema(bettershifts.Config{Resolvers: &resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
