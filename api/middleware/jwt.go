package middleware

import (
	"fmt"
	db "kangym/db"
	errors "kangym/utils/errors"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, s db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			errors.PermissionDenied(w)
			return
		}

		if !token.Valid {
			errors.PermissionDenied(w)
			return
		}

		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			errors.PermissionDenied(w)
			return
		}
		account, err := s.GetAccountById(id)
		if err != nil {
			errors.PermissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["accountId"] != strconv.Itoa(account.ID) {
			errors.PermissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func CreateJWT(accountId string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt": 15000,
		"accountId": accountId,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
