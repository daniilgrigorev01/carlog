package storage

type Storage interface {
	CreateJournal(name string, mileage int) (Journal, error)
	GetJournalById(id int) (Journal, error)
	GetAllJournals() ([]Journal, error)
	UpdateJournalName(id int, name string) (Journal, error)
	UpdateJournalMileage(id int, mileage int) (Journal, error)
	DeleteJournal(id int) error
}
