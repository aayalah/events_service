package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github/eventApp/internal/models"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const index = "events"

type EventSearchRepository struct {
	es *elasticsearch.TypedClient
}

type EventSearch struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"groupId"`
	Name        string    `json:"name"`
	Time        time.Time `json:"time"`
	Location    string    `json:"location"`
	LocationGeo GeoPoint  `json:"locationGeo"`
	DanceStyles []string  `json:"danceStyles"`
	Type        string    `json:"type"`
	Levels      []string  `json:"levels"`
}

type GeoPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func NewEventSearchRepository(es *elasticsearch.TypedClient, ctx context.Context) (*EventSearchRepository, error) {
	esr := &EventSearchRepository{es}

	err := createIndices(es, ctx)
	if err != nil {
		return nil, err
	}

	return esr, nil
}

func createIndices(es *elasticsearch.TypedClient, ctx context.Context) error {

	exists, err := es.Indices.Exists(index).IsSuccess(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id":          types.NewLongNumberProperty(),
			"groupId":     types.NewLongNumberProperty(),
			"name":        types.NewTextProperty(),
			"time":        types.NewDateProperty(),
			"location":    types.NewKeywordProperty(),
			"locationGeo": types.NewGeoPointProperty(),
			"danceStyles": types.NewKeywordProperty(),
			"type":        types.NewKeywordProperty(),
			"levels":      types.NewKeywordProperty(),
		},
	}

	_, err = es.Indices.Create(index).Mappings(mappings).Do(ctx)
	return err
}

func (s *EventSearchRepository) IndexEvent(event *models.Event, ctx context.Context) error {

	e := &EventSearch{
		Name:     event.Name,
		GroupID:  event.GroupID,
		Time:     event.Time,
		Location: event.Location,
		LocationGeo: GeoPoint{
			Latitude:  event.Latitude,
			Longitude: event.Longitude,
		},
		DanceStyles: event.DanceStyles,
		Type:        event.Type,
		Levels:      event.Levels,
	}

	_, err := s.es.Index("events").Request(e).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventSearchRepository) GetEvents(lat, long, distance float64, ctx context.Context) ([]*models.Event, error) {
	query := &types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{
					GeoDistance: &types.GeoDistanceQuery{
						Distance: fmt.Sprintf("%.2fkm", distance), // Distance in kilometers
						GeoDistanceQuery: map[string]types.GeoLocation{
							"locationGeo": types.LatLonGeoLocation{
								Lat: types.Float64(lat),
								Lon: types.Float64(long),
							},
						},
					},
				},
			},
		},
	}

	resp, err := s.es.Search().Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for _, hit := range resp.Hits.Hits {

		eventSearch := &EventSearch{}
		err := json.Unmarshal(hit.Source_, eventSearch)
		if err != nil {
			return nil, err
		}

		event := &models.Event{
			Name:        eventSearch.Name,
			GroupID:     eventSearch.GroupID,
			Time:        eventSearch.Time,
			Location:    eventSearch.Location,
			Latitude:    eventSearch.LocationGeo.Latitude,
			Longitude:   eventSearch.LocationGeo.Longitude,
			DanceStyles: eventSearch.DanceStyles,
			Type:        eventSearch.Type,
			Levels:      eventSearch.Levels,
		}
		events = append(events, event)
	}

	return events, nil
}
