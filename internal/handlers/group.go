package handlers

import (
	"context"
	"encoding/json"
	"github/eventApp/internal/service"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const groupIDParam = "groupId"

func CreateGroup(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading create group body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		group := &service.CreateGroupRequest{}

		err = json.Unmarshal(body, group)
		if err != nil {
			log.Printf("Error unmarshalling group body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		createdGroup, err := s.CreateGroup(group, ctx)
		if err != nil {
			log.Printf("Error creating group: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(createdGroup)
		if err != nil {
			log.Printf("Error marshalling created group response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}

func UpdateGroup(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		groupID := p.ByName(groupIDParam)
		ctx := context.Background()

		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading update group body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		group := &service.UpdateGroupRequest{}

		err = json.Unmarshal(body, group)
		if err != nil {
			log.Printf("Error unmarshalling group body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updatedGroup, err := s.UpdateGroup(groupIDint, group, ctx)
		if err != nil {
			log.Printf("Error updating group: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(updatedGroup)
		if err != nil {
			log.Printf("Error marshalling updated group response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)

	}
}

func GetGroups(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		ctx := context.Background()

		query := r.URL.Query()

		// Access specific parameters
		city := query.Get("city")
		country := query.Get("country")

		groups, err := s.GetGroups(city, country, ctx)
		if err != nil {
			log.Printf("Error fetching groups: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(groups)
		if err != nil {
			log.Printf("Error marshalling get groups response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}
