package handlers

import (
	"context"
	"encoding/json"
	"github/eventApp/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddUserToGroup(gtus *service.GroupToUserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.Background()

		groupID := p.ByName(groupIDParam)
		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		userID := p.ByName(userIDParam)

		userIDint, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Error converting user id param to int: %v", err)
		}

		autgr := &service.AddUserToGroupRequest{
			GroupID: groupIDint,
			UserID:  userIDint,
		}

		addedUserToGroup, err := gtus.AddUserToGroup(autgr, ctx)
		if err != nil {
			log.Printf("Error adding user to group: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(addedUserToGroup)
		if err != nil {
			log.Printf("Error marshalling adding user to group response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)

	}
}

func RemoveUserFromGroup(gtus *service.GroupToUserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.Background()

		groupID := p.ByName(groupIDParam)
		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		userID := p.ByName(userIDParam)

		userIDint, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Error converting user id param to int: %v", err)
		}

		rufgr := &service.RemoveUserFromGroupRequest{
			GroupID: groupIDint,
			UserID:  userIDint,
		}

		err = gtus.RemoveUserFromGroup(rufgr, ctx)
		if err != nil {
			log.Printf("Error removing user from group: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
