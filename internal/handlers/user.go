package handlers

import (
	"context"
	"encoding/json"
	"github/eventApp/internal/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const userIDParam = "userId"

func CreateUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading create users body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := &service.CreateUserRequest{}

		err = json.Unmarshal(body, user)
		if err != nil {
			log.Printf("Error unmarshalling users body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		createdUser, err := s.CreateUser(user, ctx)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(createdUser)
		if err != nil {
			log.Printf("Error marshalling created user response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}

func UpdateUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := p.ByName(userIDParam)
		ctx := context.Background()

		userIDint, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Error converting user id param to int: %v", err)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading update users body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := &service.UpdateUserRequest{}

		err = json.Unmarshal(body, user)
		if err != nil {
			log.Printf("Error unmarshalling users body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updatedUser, err := s.UpdateUser(userIDint, user, ctx)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(updatedUser)
		if err != nil {
			log.Printf("Error marshalling updated user response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)

	}
}

func GetUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := p.ByName(userIDParam)
		ctx := context.Background()

		userIDint, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Error converting user id param to int: %v", err)
		}

		user, err := s.GetUser(userIDint, ctx)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshalling created user response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}
