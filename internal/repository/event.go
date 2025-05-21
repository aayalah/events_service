package repository

import (
	"context"
	"github/eventApp/internal/models"

	"github.com/uptrace/bun"
)

type EventRepository struct {
	db *bun.DB
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:u"`

	ID      int64 `bun:",pk,autoincrement,nullzero"`
	GroupID int64
	Name    string
	Date    string
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
		Name:    event.Name,
		GroupID: event.GroupID,
		Date:    event.Date,
	}

	createdEvent := &Event{}

	err := s.db.NewInsert().Model(e).Returning("*").Scan(ctx, createdEvent)
	if err != nil {
		return nil, err
	}

	ce := &models.Event{
		ID:      createdEvent.ID,
		Name:    createdEvent.Name,
		GroupID: createdEvent.GroupID,
		Date:    createdEvent.Date,
	}

	return ce, nil
}

func (s *EventRepository) UpdateEvent(id int64, event *models.Event, ctx context.Context) (*models.Event, error) {

	e := &Event{
		Name:    event.Name,
		GroupID: event.GroupID,
		Date:    event.Date,
	}

	updatedEvent := &Event{}

	err := s.db.NewUpdate().Model(e).Where("id = ?", id).Returning("*").Scan(ctx, updatedEvent)
	if err != nil {
		return nil, err
	}

	ue := &models.Event{
		ID:      updatedEvent.ID,
		Name:    updatedEvent.Name,
		GroupID: updatedEvent.GroupID,
		Date:    updatedEvent.Date,
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
			ID:      e.ID,
			Name:    e.Name,
			GroupID: e.GroupID,
			Date:    e.Date,
		})
	}

	return mgs, nil
}
