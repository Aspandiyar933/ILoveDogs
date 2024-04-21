package tasks

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Aspandiyar933/Ilovedogs/auth"
	"github.com/Aspandiyar933/Ilovedogs/store"
	"github.com/Aspandiyar933/Ilovedogs/typeslink"
	"github.com/Aspandiyar933/Ilovedogs/utils"
	"github.com/gorilla/mux"
)

var (
	errContentRequired    = errors.New("content is required")
	errMonthRequired      = errors.New("month is required")
	errBreedRequired      = errors.New("breed is required")
	errGenderRequired     = errors.New("gender is required")
	errVaccinatedRequired = errors.New("vaccinated is required")
)

type PostService struct {
	store store.Store
}

func NewPostService(s store.Store) *PostService {
	return &PostService{
		store: s,
	}
}

func (s *PostService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/dog", auth.WithJWTAuth(s.HandleCreateDog, s.store)).Methods("POST")
	r.HandleFunc("/dog/{id}", auth.WithJWTAuth(s.HandleGetDog, s.store)).Methods("GET")
}

func (s *PostService) HandleCreateDog(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var post *typeslink.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Invalid request payload!"}) 
		return
	}
	if err := validatePostPayload(post); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "Invalid request payload!"}) 
		return
	}

	p, err := s.store.CreatePost(post)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: err.Error()}) 
		return
	}
	utils.WriteJSON(w, http.StatusCreated, p)

}

func (s *PostService) HandleGetDog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, typeslink.ErrorResponse{Error: "Id is required"})
		return
	}
	p, err := s.store.GetPost(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, typeslink.ErrorResponse{Error: "post not found"})
		return 
	}
	utils.WriteJSON(w, http.StatusOK, p)
}

func validatePostPayload(post *typeslink.Post) error {
	switch {
	case post.Content == "":
		return errContentRequired
	case post.Month == 0:
		return errMonthRequired
	case post.Breed == "":
		return errBreedRequired
	case post.Gender == "":
		return errGenderRequired
	case post.Vaccinated == "":
		return errVaccinatedRequired
	default:
		return nil
	}
}
