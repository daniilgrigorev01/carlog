package storage

import "time"

type Journal struct {
	ID                  int
	Name                string
	CurrentMileage      int
	MileageLastUpdateAt time.Time
}

type Record struct {
	ID          int
	JournalID   int
	Name        string
	SKU         string
	ReplacedAt  time.Time
	NextDate    *time.Time
	Mileage     int
	NextMileage *int
	CreatedAt   time.Time
}

type UpdateRecordParams struct {
	Name        *string
	SKU         *string
	ReplacedAt  *time.Time
	NextDate    *time.Time
	Mileage     *int
	NextMileage *int
}
