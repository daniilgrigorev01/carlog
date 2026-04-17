package storage

import "time"

type Journal struct {
	ID                  int64
	Name                string
	CurrentMileage      int
	MileageLastUpdateAt time.Time
}

type Record struct {
	ID          int64
	JournalID   int64
	Name        string
	SKU         string
	ReplacedAt  time.Time
	NextDate    *time.Time
	Mileage     int
	NextMileage *int
	CreatedAt   time.Time
}

type CreateRecordParams struct {
	JournalID   int64
	Name        string
	SKU         string
	ReplacedAt  time.Time
	NextDate    *time.Time
	Mileage     int
	NextMileage *int
}

type UpdateRecordParams struct {
	Name        *string
	SKU         *string
	ReplacedAt  *time.Time
	NextDate    *time.Time
	Mileage     *int
	NextMileage *int
}
