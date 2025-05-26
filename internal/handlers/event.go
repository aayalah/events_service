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

const eventIDParam = "eventId"

func CreateEvent(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func UpdateEvent(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func GetEvents(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func GetEventsByDistance(s *service.EventService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		ctx := context.Background()

		lat := r.URL.Query().Get("lat")
		latFloat64, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			log.Printf("Error converting latitude to float64: %v", err)
			return
		}

		long := r.URL.Query().Get("long")
		longFloat64, err := strconv.ParseFloat(long, 64)
		if err != nil {
			log.Printf("Error converting longitude to float64: %v", err)
			return
		}

		distance := r.URL.Query().Get("distance")
		distanceFloat64, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			log.Printf("Error converting latitude to float64: %v", err)
			return
		}

		events, err := s.GetEventsByDistance(latFloat64, longFloat64, distanceFloat64, ctx)
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
