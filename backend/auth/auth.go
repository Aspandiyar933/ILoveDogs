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
	return func(w http.ResponseWriter, r *http.Request)  {
		tokenString := GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Panicln("Failed to authenticate token")
			PermissionDenied(w)
			return
		}
		if !token.Valid {
			log.Panicln("Failed to authenticate token")
			PermissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		id := claims["ID"].(string)

		_, err = store.GetUserByID(id)
		if err != nil {
			log.Panicln("Failed to get user id")
			PermissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w,http.StatusUnauthorized, typeslink.ErrorResponse{
		Error: fmt.Errorf("Permission denied").Error(),
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}
	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(p string) (*jwt.Token, error) {
	secret := config.Envs.JWTSecret

	return jwt.Parse(p, func(p *jwt.Token) (interface{}, error){
		if _, ok := p.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpexted signing method: %v", p.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func HashPasswords(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil 
}

func CreateJWT(secret []byte, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(userId)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}