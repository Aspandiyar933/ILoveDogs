package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Aspandiyar933/Ilovedogs/auth"
	"github.com/Aspandiyar933/Ilovedogs/config"
	"github.com/Aspandiyar933/Ilovedogs/store"
	"github.com/Aspandiyar933/Ilovedogs/typeslink"
	"github.com/Aspandiyar933/Ilovedogs/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	errNameRequired       = errors.New("name is required")
	errUserNameRequired   = errors.New("username is required")
	errEmailRequired      = errors.New("email is required")
	errPasswordRequired   = errors.New("password is required")
	errVaccinatedRequired = errors.New("vaccinated is required")
)

type UserService struct {
	store store.Store
}

func NewUserService(s store.Store) *UserService {
	return &UserService{
		store: s,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.HandleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
}

func (s *UserService) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    var login typeslink.LoginRequest
    if err := json.Unmarshal(body, &login); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    foundUser, err := s.store.GetUserByEmail(login.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password)); err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    token, err := createAndSetAuthCookieL(foundUser.Email, w)
    if err != nil {
        http.Error(w, "Error creating session", http.StatusInternalServerError)
        return
    }

    // Respond with token
    utils.WriteJSON(w, http.StatusOK, token)
}


// Assuming createAndSetAuthCookie generates and sets the JWT token
func createAndSetAuthCookieL(userID string, w http.ResponseWriter) (string, error) {
	// Generate JWT token with user ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	// Sign the token with a secret
	secret := []byte(config.Envs.JWTSecret)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Example: token expires in 24 hours
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return tokenString, nil
}


func (s *UserService) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload *typeslink.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, typeslink.ErrorResponse{Error: "Invalid request payload!"})
		return
	}

	if err = validateUserPayload(payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, typeslink.ErrorResponse{Error: err.Error()})
		return
	}

	hashedPW, err := auth.HashPasswords(payload.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Error creating user"})
		return
	}
	payload.Password = hashedPW

	u, err := s.store.CreateUser(payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Error creating user"})
		return
	}
	token, err :=  createAndSetAuthCoolie(u.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Error creating session"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)
}

func validateUserPayload(user *typeslink.User) error {
	switch {
	case user.Name == "":
		return errNameRequired
	case user.UserName == "":
		return errUserNameRequired
	case user.Email == "":
		return errEmailRequired
	case user.Password == "":
		return errPasswordRequired
	default:
		return nil
	}
}

func validateLoginRequestPayload(login *typeslink.LoginRequest) error {
	switch {
	case login.Email == "":
		return errNameRequired
	case login.Password == "":
		return errEmailRequired
	default:
		return nil
	}
}

func createAndSetAuthCoolie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, id)
	if err != nil {
		return "", nil 
	}

	http.SetCookie(w, &http.Cookie{
		Name: "Authorization",
		Value: token,
	})

	return token, nil
}