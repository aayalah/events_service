package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github/eventApp/config"
	"github/eventApp/internal/repository"
	"github/eventApp/internal/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {

	config, err := config.New()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
		return
	}

	dsn := "postgres://postgres:postgres@localhost:5455/postgres?sslmode=disable"
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err = sqlDb.Ping()
	if err != nil {
		log.Fatalf("Error connecting to postgres DB: %v", err)
	}

	db := bun.NewDB(sqlDb, pgdialect.New())

	/* repositories */
	groupToUserRep, err := repository.NewGroupToUserRepository(db, context.Background())
	if err != nil {
		log.Fatalf("Error creating group to user repository: %v", err)
	}

	userRep, err := repository.NewUserRepository(db, context.Background())
	if err != nil {
		log.Fatalf("Error creating user repository: %v", err)
	}

	groupRep, err := repository.NewGroupRepository(db, context.Background())
	if err != nil {
		log.Fatalf("Error creating group repository: %v", err)
	}

	eventRep, err := repository.NewEventRepository(db, context.Background())
	if err != nil {
		log.Fatalf("Error creating event repository: %v", err)
	}

	userService := service.NewUserService(userRep)
	groupService := service.NewGroupService(groupRep)
	eventService := service.NewEventService(eventRep)
	groupToUserService := service.NewGroupToUserService(groupToUserRep)

	/*server
	 */

	router := httprouter.New()
	router.POST("/users", createUser(userService))
	router.GET("/users/:userId", getUser(userService))
	router.PUT("/users/:userId", updateUser(userService))

	router.POST("/groups", createGroup(groupService))
	router.GET("/groups", getGroups(groupService))
	router.PUT("/groups/:groupId", updateGroup(groupService))
	router.GET("/groups/:groupId/events", getEvents(eventService))
	router.POST("/groups/:groupId/events", createEvent(eventService))
	router.PUT("/groups/:groupId/events/:eventId", updateEvent(eventService))

	router.POST("/groups/:groupId/users/:userId", addUserToGroup(groupToUserService))
	router.DELETE("/groups/:groupId/users/:userId", removeUserFromGroup(groupToUserService))

	http.ListenAndServe(fmt.Sprintf(":%v", config.PORT), router)

}

const userIDParam = "userId"
const groupIDParam = "groupId"
const eventIDParam = "eventId"

/* handlers */

func createUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		body, err := ioutil.ReadAll(r.Body)
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

func updateUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func getUser(s *service.UserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func createGroup(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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

func updateGroup(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func getGroups(s *service.GroupService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func createEvent(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		groupID := p.ByName(groupIDParam)
		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading create event body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		event := &service.CreateEventRequest{}

		err = json.Unmarshal(body, event)
		if err != nil {
			log.Printf("Error unmarshalling event body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		event.GroupID = groupIDint

		ctx := context.Background()

		createdEvent, err := s.CreateEvent(event, ctx)
		if err != nil {
			log.Printf("Error creating event: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(createdEvent)
		if err != nil {
			log.Printf("Error marshalling created event response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}

func updateEvent(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		groupID := p.ByName(groupIDParam)
		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		eventID := p.ByName(eventIDParam)
		ctx := context.Background()

		eventIDint, err := strconv.ParseInt(eventID, 10, 64)
		if err != nil {
			log.Printf("Error converting event id param to int: %v", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading update event body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		event := &service.UpdateEventRequest{}

		err = json.Unmarshal(body, event)
		if err != nil {
			log.Printf("Error unmarshalling event body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		event.GroupID = groupIDint

		updatedEvent, err := s.UpdateEvent(eventIDint, event, ctx)
		if err != nil {
			log.Printf("Error updating event: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(updatedEvent)
		if err != nil {
			log.Printf("Error marshalling updated event response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)

	}
}

func getEvents(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		ctx := context.Background()

		groupID := p.ByName(groupIDParam)
		groupIDint, err := strconv.ParseInt(groupID, 10, 64)
		if err != nil {
			log.Printf("Error converting group id param to int: %v", err)
		}

		events, err := s.GetEvents(groupIDint, ctx)
		if err != nil {
			log.Printf("Error fetching events: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(events)
		if err != nil {
			log.Printf("Error marshalling get events response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}

func addUserToGroup(gtus *service.GroupToUserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func removeUserFromGroup(gtus *service.GroupToUserService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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
