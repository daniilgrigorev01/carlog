package storage

type Storage interface {
	CreateJournal(name string, mileage int) (Journal, error)
	GetJournalByID(id int) (Journal, error)
	GetAllJournals() ([]Journal, error)
	UpdateJournalName(id int, name string) (Journal, error)
	UpdateJournalMileage(id int, mileage int) (Journal, error)
	DeleteJournal(id int) error

	AddRecord(record Record) (Record, error)
	GetRecordByID(id int) (Record, error)
	GetAllRecordsByJournalID(journalID int) ([]Record, error)
	UpdateRecord(id int, params UpdateRecordParams) (Record, error)
	DeleteRecord(id int) error
}
