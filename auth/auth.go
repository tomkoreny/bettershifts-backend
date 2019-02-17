package auth

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/lordpuma/bettershifts/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var userCtxKey = &contextKey{"user"}
type contextKey struct {
	name string
}

func validateAndGetUserID(c string, db *gorm.DB) (string, error) {
  var token models.Token
  db.First(&token, models.Token{Token: c})
  return token.UserID, nil;
}

func getUserByID(id string, db *gorm.DB) models.User {
  var user models.User
  db.First(&user, models.User{ID: id})
  return user; 
}

func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				next.ServeHTTP(w, r)
				return
			}


			userId, err := validateAndGetUserID(token[7:len(token)], db)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			user := getUserByID(userId, db)

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {
	raw := ctx.Value(userCtxKey)
  if raw == nil {
    return nil
  }
  ret := raw.(models.User)
	return &ret
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
    // Since we'll be getting the hashed password from the DB it
    // will be a string so we'll need to convert it to a byte slice
    byteHash := []byte(hashedPwd)
    plainbyteHash := []byte(plainPwd)
    err := bcrypt.CompareHashAndPassword(byteHash, plainbyteHash)
    if err != nil {
        return false
    }
    
    return true
  }

  func HashAndSalt(pwd string) string {
    
    // Use GenerateFromPassword to hash & salt pwd.
    // MinCost is just an integer constant provided by the bcrypt
    // package along with DefaultCost & MaxCost. 
    // The cost can be any value you want provided it isn't lower
    // than the MinCost (4)
    pwdByte := []byte(pwd)
    hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    // GenerateFromPassword returns a byte slice so we need to
    // convert the bytes to a string and return it
    return string(hash)
  }
