package models

import "time"

type Event struct {
	ID          int64
	Name        string
	Time        time.Time
	Latitude    float64
	Longitude   float64
	Location    string
	GroupID     int64
	DanceStyles []string
	Type        string
	Levels      []string
}
