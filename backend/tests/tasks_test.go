package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aspandiyar933/Ilovedogs/tasks"
	"github.com/Aspandiyar933/Ilovedogs/typeslink"
	"github.com/gorilla/mux"
)

func TestCreatePost(t *testing.T) {
	ms := &MockStore{}
	service := tasks.NewPostService(ms)
	t.Run("should create the post", func(t *testing.T) {
		payload := &typeslink.Post{
			Content: "",
			Month: 5,
			Breed: "african",
			Gender: "F",
			Vaccinated: "Yes",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/dog", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/dog", service.HandleCreateDog)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expacted status code %d, god %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetPost(t *testing.T) {
	ms := &MockStore{}
	service := tasks.NewPostService(ms)

	t.Run("should return the post", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/god/41", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/dog/{id}", service.HandleGetDog)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Error("invalid status code, it should fail")
		}
	})
}