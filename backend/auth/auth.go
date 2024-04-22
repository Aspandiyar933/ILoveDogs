package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Aspandiyar933/Ilovedogs/config"
	"github.com/Aspandiyar933/Ilovedogs/store"
	"github.com/Aspandiyar933/Ilovedogs/typeslink"
	"github.com/Aspandiyar933/Ilovedogs/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			log.Printf("Failed to authenticate token: %v", err)
			PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["userID"].(string)
		if !ok {
			log.Println("UserID not found in token claims")
			PermissionDenied(w)
			return
		}

		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Failed to parse user ID: %v", err)
			PermissionDenied(w)
			return
		}

		_, err = store.GetUserByID(strconv.FormatInt(id, 10))
		if err != nil {
			log.Printf("Failed to get user by ID: %v", err)
			PermissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, typeslink.ErrorResponse{
		Error: "Permission denied",
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := []byte(config.Envs.JWTSecret)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CreateJWT(userID int64) (string, error) {
	secret := []byte(config.Envs.JWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.FormatInt(userID, 10),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
