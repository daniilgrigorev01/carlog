package storage

import (
	"time"
)

type MockStorage struct {
	journals      []Journal
	records       []Record
	nextJournalID int
	nextRecordID  int
}

func (m *MockStorage) CreateJournal(name string, mileage int) (Journal, error) {
	newJournal := Journal{
		ID:                  m.nextJournalID,
		Name:                name,
		CurrentMileage:      mileage,
		MileageLastUpdateAt: time.Now(),
	}

	m.journals = append(m.journals, newJournal)
	m.nextJournalID++

	return newJournal, nil
}

func (m *MockStorage) GetJournalByID(id int) (Journal, error) {
	for i := range m.journals {
		if m.journals[i].ID == id {
			return m.journals[i], nil
		}
	}

	return Journal{}, ErrNotFound
}

func (m *MockStorage) GetAllJournals() ([]Journal, error) {
	return m.journals, nil
}

func (m *MockStorage) UpdateJournalName(id int, name string) (Journal, error) {
	for i := range m.journals {
		if m.journals[i].ID == id {
			m.journals[i].Name = name
			return m.journals[i], nil
		}
	}

	return Journal{}, ErrNotFound
}

func (m *MockStorage) UpdateJournalMileage(id int, mileage int) (Journal, error) {
	for i := range m.journals {
		if m.journals[i].ID == id {
			m.journals[i].CurrentMileage = mileage
			m.journals[i].MileageLastUpdateAt = time.Now()

			return m.journals[i], nil
		}
	}

	return Journal{}, ErrNotFound
}

func (m *MockStorage) DeleteJournal(id int) error {
	for i := range m.journals {
		if m.journals[i].ID == id {
			m.journals = append(m.journals[:i], m.journals[i+1:]...)

			return nil
		}
	}

	return ErrNotFound
}

func (m *MockStorage) AddRecord(record Record) (Record, error) {
	record.ID = m.nextRecordID

	m.records = append(m.records, record)
	m.nextRecordID++

	return record, nil
}

func (m *MockStorage) GetRecordByID(id int) (Record, error) {
	for i := range m.records {
		if m.records[i].ID == id {
			return m.records[i], nil
		}
	}

	return Record{}, ErrNotFound
}

func (m *MockStorage) GetAllRecordsByJournalID(journalID int) ([]Record, error) {
	var result []Record

	for i := range m.records {
		if m.records[i].JournalID == journalID {
			result = append(result, m.records[i])
		}
	}

	return result, nil
}

func (m *MockStorage) UpdateRecord(id int, params UpdateRecordParams) (Record, error) {
	for i := range m.records {
		if m.records[i].ID == id {
			if params.Name != nil {
				m.records[i].Name = *params.Name
			}
			if params.SKU != nil {
				m.records[i].SKU = *params.SKU
			}
			if params.ReplacedAt != nil {
				m.records[i].ReplacedAt = *params.ReplacedAt
			}
			if params.Mileage != nil {
				m.records[i].Mileage = *params.Mileage
			}

			return m.records[i], nil
		}
	}

	return Record{}, ErrNotFound
}

func (m *MockStorage) DeleteRecord(id int) error {
	for i := range m.records {
		if m.records[i].ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)

			return nil
		}
	}

	return ErrNotFound
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		nextJournalID: 1,
		nextRecordID:  1,
	}
}
