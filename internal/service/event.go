package service

import (
	"context"
	"github/eventApp/internal/models"
)

type eventRep interface {
	CreateEvent(event *models.Event, ctx context.Context) (*models.Event, error)
	UpdateEvent(id int64, event *models.Event, ctx context.Context) (*models.Event, error)
	GetEvents(groupID int64, ctx context.Context) ([]*models.Event, error)
}

type EventService struct {
	eventRep eventRep
}

func NewEventService(eventRep eventRep) *EventService {
	return &EventService{
		eventRep,
	}
}

type CreateEventRequest struct {
	Name    string `json:"name"`
	GroupID int64  `json:"groupId"`
	Date    string `json:"date"`
}

type CreateEventResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	GroupID int64  `json:"groupId"`
	Date    string `json:"date"`
}

func (e *EventService) CreateEvent(cer *CreateEventRequest, ctx context.Context) (*CreateEventResponse, error) {

	event := &models.Event{
		Name:    cer.Name,
		GroupID: cer.GroupID,
		Date:    cer.Date,
	}

	createdEvent, err := e.eventRep.CreateEvent(event, ctx)
	if err != nil {
		return nil, err
	}

	ceResp := &CreateEventResponse{
		ID:      createdEvent.ID,
		Name:    createdEvent.Name,
		GroupID: createdEvent.GroupID,
		Date:    createdEvent.Date,
	}

	return ceResp, nil

}

type GetEventResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	GroupID int64  `json:"groupId"`
	Date    string `json:"date"`
}

func (s *EventService) GetEvents(groupID int64, ctx context.Context) ([]*GetEventResponse, error) {

	events, err := s.eventRep.GetEvents(groupID, ctx)
	if err != nil {
		return nil, err
	}

	eventsResp := make([]*GetEventResponse, 0, len(events))

	for _, e := range events {
		eventsResp = append(eventsResp, &GetEventResponse{
			ID:      e.ID,
			Name:    e.Name,
			GroupID: e.GroupID,
			Date:    e.Date,
		})
	}

	return eventsResp, nil
}

type UpdateEventRequest struct {
	Name    string `json:"name"`
	GroupID int64  `json:"groupId"`
	Date    string `json:"date"`
}

type UpdateEventResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	GroupID int64  `json:"groupId"`
	Date    string `json:"date"`
}

func (s *EventService) UpdateEvent(id int64, uer *UpdateEventRequest, ctx context.Context) (*UpdateEventResponse, error) {

	event := &models.Event{
		Name:    uer.Name,
		GroupID: uer.GroupID,
		Date:    uer.Date,
	}

	updatedEvent, err := s.eventRep.UpdateEvent(id, event, ctx)
	if err != nil {
		return nil, err
	}

	ueResp := &UpdateEventResponse{
		ID:      updatedEvent.ID,
		Name:    updatedEvent.Name,
		GroupID: updatedEvent.GroupID,
		Date:    updatedEvent.Date,
	}

	return ueResp, nil

}
