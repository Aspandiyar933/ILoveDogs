package users

import (
	"encoding/json"
	"errors"
	"go/token"
	"io"
	"net/http"

	"github.com/Aspandiyar933/Ilovedogs/auth"
	"github.com/Aspandiyar933/Ilovedogs/config"
	"github.com/Aspandiyar933/Ilovedogs/store"
	"github.com/Aspandiyar933/Ilovedogs/typeslink"
	"github.com/Aspandiyar933/Ilovedogs/utils"
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

	var user *typeslink.User
	foundUser, err := s.store.GetUserByEmail(user.Email)
    if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, typeslink.ErrorResponse{Error: "User not found"})
        return
    }
	var login *typeslink.LoginRequest
	err = json.Unmarshal(body, &login)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, typeslink.ErrorResponse{Error: "Invalid request payload!"})
		return
	}

	if err = validateLoginRequestPayload(login); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, typeslink.ErrorResponse{Error: err.Error()})
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password)); err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, typeslink.ErrorResponse{Error: "Invalid password"})
		return
	}

	token, err :=  createAndSetAuthCoolie(, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Error creating session"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)
	
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