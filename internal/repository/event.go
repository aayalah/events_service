package repository

import (
	"context"
	"github/eventApp/internal/models"
	"time"

	"github.com/uptrace/bun"
)

type EventRepository struct {
	db *bun.DB
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:u"`

	ID          int64     `bun:",pk,autoincrement,nullzero"`
	GroupID     int64     `bun:",notnull"`
	Name        string    `bun:",notnull"`
	Time        time.Time `bun:"time,notnull"`
	Location    string    `bun:",notnull"`
	Latitude    float64   `bun:",notnull"`
	Longitude   float64   `bun:",notnull"`
	DanceStyles []string
	Type        string
	Levels      []string
}

func NewEventRepository(db *bun.DB, ctx context.Context) (*EventRepository, error) {
	usr := &EventRepository{db}
	err := usr.createEventTable(ctx)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *EventRepository) createEventTable(ctx context.Context) error {
	_, err := s.db.NewCreateTable().IfNotExists().Model((*Event)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventRepository) CreateEvent(event *models.Event, ctx context.Context) (*models.Event, error) {

	e := &Event{
		Name:        event.Name,
		GroupID:     event.GroupID,
		Time:        event.Time,
		Latitude:    event.Latitude,
		Longitude:   event.Longitude,
		Location:    event.Location,
		DanceStyles: event.DanceStyles,
		Type:        event.Type,
		Levels:      event.Levels,
	}

	createdEvent := &Event{}

	err := s.db.NewInsert().Model(e).Returning("*").Scan(ctx, createdEvent)
	if err != nil {
		return nil, err
	}

	ce := &models.Event{
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

	return ce, nil
}

func (s *EventRepository) UpdateEvent(id int64, event *models.Event, ctx context.Context) (*models.Event, error) {

	e := &Event{
		Name:        event.Name,
		GroupID:     event.GroupID,
		Time:        event.Time,
		Latitude:    event.Latitude,
		Longitude:   event.Longitude,
		Location:    event.Location,
		DanceStyles: event.DanceStyles,
		Type:        event.Type,
		Levels:      event.Levels,
	}

	updatedEvent := &Event{}

	err := s.db.NewUpdate().Model(e).Where("id = ?", id).Returning("*").Scan(ctx, updatedEvent)
	if err != nil {
		return nil, err
	}

	ue := &models.Event{
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

	return ue, nil
}

func (s *EventRepository) GetEvents(groupID int64, ctx context.Context) ([]*models.Event, error) {
	var events []Event

	err := s.db.NewSelect().Model(&events).Where("group_id = ?", groupID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	mgs := make([]*models.Event, 0, len(events))

	for _, e := range events {
		mgs = append(mgs, &models.Event{
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

	return mgs, nil
}
