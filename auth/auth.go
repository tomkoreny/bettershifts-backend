package auth

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/lordpuma/bettershifts-backend/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
  "errors"
)

var userCtxKey = &contextKey{"user"}
type contextKey struct {
	name string
}

func validateAndGetUserID(c string, db *gorm.DB) (string, error) {
  var token models.Token
  db.Where("token = ?", c).First(&token)
  if (token.ID == "") {
    return "", errors.New("invalid-token")
  }
  return token.UserID, nil;
}

func getUserByID(id string, db *gorm.DB) models.User {
  var user models.User
  db.Where("id", id).First(&user)
  return user; 
}

func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" || len(token) != 39 {
				next.ServeHTTP(w, r)
				return
			} 

			userId, err := validateAndGetUserID(token[7:len(token)], db)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user := getUserByID(userId, db)

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *models.User {
	raw := ctx.Value(userCtxKey)
  if raw == nil {
    return nil
  }
  ret := raw.(models.User)
	return &ret
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
    byteHash := []byte(hashedPwd)
    plainbyteHash := []byte(plainPwd)
    err := bcrypt.CompareHashAndPassword(byteHash, plainbyteHash)
    if err != nil {
        return false
    } 
    return true
  }

  func HashAndSalt(pwd string) string { 
    pwdByte := []byte(pwd)
    hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    return string(hash)
  }
