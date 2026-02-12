package storage

import "time"

type Log struct {
	ID                  int
	Name                string
	CurrentMileage      int
	MileageLastUpdateAt time.Time
}

type Record struct {
	ID          int
	LogID       int
	Name        string
	Article     string
	ReplacedAt  time.Time
	NextDate    *time.Time
	Mileage     int
	NextMileage *int
}
