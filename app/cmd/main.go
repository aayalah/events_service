package main

import (
	"context"
	"database/sql"
	"fmt"
	"github/eventApp/config"
	"github/eventApp/internal/handlers"
	"github/eventApp/internal/middleware"
	"github/eventApp/internal/repository"
	"github/eventApp/internal/service"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
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
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: config.ELASTIC_SEARCH_ADDRESSES,
	})
	if err != nil {
		log.Fatalf("Error creating the elastic search client: %v", err)
	}

	_, err = es.Info().Do(context.Background())
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch: %v", err)
	}

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

	eventSearchRep, err := repository.NewEventSearchRepository(es, context.Background())
	if err != nil {
		log.Fatalf("Error creating event search repository: %v", err)
	}

	userService := service.NewUserService(userRep)
	groupService := service.NewGroupService(groupRep)
	eventService := service.NewEventService(eventRep, eventSearchRep)
	groupToUserService := service.NewGroupToUserService(groupToUserRep)
	loginService := service.NewLoginService(config.JWTSECRET, userRep)

	/*server
	 */

	router := httprouter.New()
	router.POST("/login", handlers.Login(loginService))

	router.POST("/users", handlers.CreateUser(userService))
	router.GET("/users/:userId", middleware.Auth(config.JWTSECRET, handlers.GetUser(userService)))
	router.PATCH("/users/:userId", middleware.Auth(config.JWTSECRET, handlers.UpdateUser(userService)))

	router.POST("/groups", middleware.Auth(config.JWTSECRET, handlers.CreateGroup(groupService)))
	router.GET("/groups", middleware.Auth(config.JWTSECRET, handlers.GetGroups(groupService)))
	router.PUT("/groups/:groupId", middleware.Auth(config.JWTSECRET, handlers.UpdateGroup(groupService)))
	router.GET("/groups/:groupId/events", middleware.Auth(config.JWTSECRET, handlers.GetEvents(eventService)))
	router.POST("/groups/:groupId/events", middleware.Auth(config.JWTSECRET, handlers.CreateEvent(eventService)))
	router.PUT("/groups/:groupId/events/:eventId", middleware.Auth(config.JWTSECRET, handlers.UpdateEvent(eventService)))

	router.POST("/groups/:groupId/users/:userId", middleware.Auth(config.JWTSECRET, handlers.AddUserToGroup(groupToUserService)))
	router.DELETE("/groups/:groupId/users/:userId", middleware.Auth(config.JWTSECRET, handlers.RemoveUserFromGroup(groupToUserService)))

	router.GET("/events", middleware.Auth(config.JWTSECRET, handlers.GetEventsByDistance(eventService)))
	http.ListenAndServe(fmt.Sprintf(":%v", config.PORT), router)

}
