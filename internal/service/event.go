package service

import (
	"context"
	"github/eventApp/internal/models"
	"log"
	"time"
)

type eventRep interface {
	CreateEvent(event *models.Event, ctx context.Context) (*models.Event, error)
	UpdateEvent(id int64, event *models.Event, ctx context.Context) (*models.Event, error)
	GetEvents(groupID int64, ctx context.Context) ([]*models.Event, error)
}

type eventSearchRep interface {
	GetEvents(lat, long, distance float64, ctx context.Context) ([]*models.Event, error)
	IndexEvent(event *models.Event, ctx context.Context) error
}

type EventService struct {
	eventRep      eventRep
	eventSearcher eventSearchRep
}

func NewEventService(eventRep eventRep, eventSearchRep eventSearchRep) *EventService {
	return &EventService{
		eventRep,
		eventSearchRep,
	}
}

type CreateEventRequest struct {
	Name        string    `json:"name"`
	GroupID     int64     `json:"groupId"`
	Time        time.Time `json:"time"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Location    string    `json:"location"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

type CreateEventResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	GroupID     int64     `json:"groupId"`
	Time        time.Time `json:"time"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Location    string    `json:"location"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

func (e *EventService) CreateEvent(cer *CreateEventRequest, ctx context.Context) (*CreateEventResponse, error) {

	event := &models.Event{
		Name:        cer.Name,
		GroupID:     cer.GroupID,
		Time:        cer.Time,
		Latitude:    cer.Latitude,
		Longitude:   cer.Longitude,
		Location:    cer.Location,
		DanceStyles: cer.DanceStyles,
		Type:        cer.Type,
		Levels:      cer.Levels,
	}

	createdEvent, err := e.eventRep.CreateEvent(event, ctx)
	if err != nil {
		return nil, err
	}

	err = e.eventSearcher.IndexEvent(createdEvent, ctx)
	if err != nil {
		log.Printf("error adding event to elastic search: %v", err)
	}

	ceResp := &CreateEventResponse{
		ID:          createdEvent.ID,
		Name:        createdEvent.Name,
		GroupID:     createdEvent.GroupID,
		Time:        createdEvent.Time,
		Latitude:    createdEvent.Latitude,
		Longitude:   createdEvent.Longitude,
		Location:    createdEvent.Location,
		DanceStyles: createdEvent.DanceStyles,
		Type:        createdEvent.Type,
		Levels:      createdEvent.Levels,
	}

	return ceResp, nil

}

type GetEventResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	GroupID     int64     `json:"groupId"`
	Time        time.Time `json:"time"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Location    string    `json:"location"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

func (e *EventService) GetEvents(groupID int64, ctx context.Context) ([]*GetEventResponse, error) {
	events, err := e.eventRep.GetEvents(groupID, ctx)
	if err != nil {
		return nil, err
	}

	eventsResp := make([]*GetEventResponse, 0, len(events))

	for _, e := range events {
		eventsResp = append(eventsResp, &GetEventResponse{
			ID:          e.ID,
			Name:        e.Name,
			GroupID:     e.GroupID,
			Time:        e.Time,
			Latitude:    e.Latitude,
			Longitude:   e.Longitude,
			Location:    e.Location,
			DanceStyles: e.DanceStyles,
			Type:        e.Type,
			Levels:      e.Levels,
		})
	}

	return eventsResp, nil
}

type UpdateEventRequest struct {
	Name        string    `json:"name"`
	GroupID     int64     `json:"groupId"`
	Time        time.Time `json:"time"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Location    string    `json:"location"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

type UpdateEventResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	GroupID     int64     `json:"groupId"`
	Time        time.Time `json:"time"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Location    string    `json:"location"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

func (e *EventService) UpdateEvent(id int64, uer *UpdateEventRequest, ctx context.Context) (*UpdateEventResponse, error) {

	event := &models.Event{
		Name:        uer.Name,
		GroupID:     uer.GroupID,
		Time:        uer.Time,
		Latitude:    uer.Latitude,
		Longitude:   uer.Longitude,
		Location:    uer.Location,
		DanceStyles: uer.DanceStyles,
		Type:        uer.Type,
		Levels:      uer.Levels,
	}

	updatedEvent, err := e.eventRep.UpdateEvent(id, event, ctx)
	if err != nil {
		return nil, err
	}

	err = e.eventSearcher.IndexEvent(updatedEvent, ctx)
	if err != nil {
		log.Printf("error adding event to elastic search: %v", err)
	}

	ueResp := &UpdateEventResponse{
		ID:          updatedEvent.ID,
		Name:        updatedEvent.Name,
		GroupID:     updatedEvent.GroupID,
		Time:        updatedEvent.Time,
		Latitude:    updatedEvent.Latitude,
		Longitude:   updatedEvent.Longitude,
		Location:    updatedEvent.Location,
		DanceStyles: updatedEvent.DanceStyles,
		Type:        updatedEvent.Type,
		Levels:      updatedEvent.Levels,
	}

	return ueResp, nil

}

func (e *EventService) GetEventsByDistance(lat, long, distance float64, ctx context.Context) ([]*GetEventResponse, error) {
	events, err := e.eventSearcher.GetEvents(lat, long, distance, ctx)
	if err != nil {
		return nil, err
	}

	eventsResp := make([]*GetEventResponse, 0, len(events))

	for _, e := range events {
		eventsResp = append(eventsResp, &GetEventResponse{
			ID:          e.ID,
			Name:        e.Name,
			GroupID:     e.GroupID,
			Time:        e.Time,
			Latitude:    e.Latitude,
			Longitude:   e.Longitude,
			Location:    e.Location,
			DanceStyles: e.DanceStyles,
			Type:        e.Type,
			Levels:      e.Levels,
		})
	}

	return eventsResp, nil
}
