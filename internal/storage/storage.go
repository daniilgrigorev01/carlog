package storage

type Storage interface {
	CreateJournal(name string, mileage int) (Journal, error)
	GetJournalByID(id int64) (Journal, error)
	GetAllJournals() ([]Journal, error)
	UpdateJournalName(id int64, name string) (Journal, error)
	UpdateJournalMileage(id int64, mileage int) (Journal, error)
	DeleteJournal(id int64) error

	AddRecord(record Record) (Record, error)
	GetRecordByID(id int64) (Record, error)
	GetAllRecordsByJournalID(journalID int64) ([]Record, error)
	UpdateRecord(id int64, params UpdateRecordParams) (Record, error)
	DeleteRecord(id int64) error

	Close() error
}
