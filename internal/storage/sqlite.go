package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type SQLiteStorage struct {
	db *sql.DB
}

func (s *SQLiteStorage) migrate() error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("storage: get migration files: %w", err)
	}

	var names []string

	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY AUTOINCREMENT,
			applied_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		)`)
	if err != nil {
		return fmt.Errorf("storage: create migration table: %w", err)
	}

	var version int
	if err = s.db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&version); err != nil {
		return fmt.Errorf("storage: get migration version: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("storage: start migration transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, name := range names[version:] {
		content, err := fs.ReadFile(migrationsFS, "migrations/"+name)
		if err != nil {
			return fmt.Errorf("storage: reading migration: %w", err)
		}

		if _, err = tx.Exec(string(content)); err != nil {
			return fmt.Errorf("storage: apply migration: %w", err)
		}

		if _, err = tx.Exec("INSERT INTO schema_migrations DEFAULT VALUES"); err != nil {
			return fmt.Errorf("storage: update migration version: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: migration failed: %w", err)
	}

	return nil
}

func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", path)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("storage: database unavailable: %w", err)
	}

	s := &SQLiteStorage{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStorage) CreateJournal(name string, mileage int) (Journal, error) {
	result, err := s.db.Exec(
		"INSERT INTO journals (name, current_mileage) VALUES (?, ?)",
		name, mileage,
	)
	if err != nil {
		return Journal{}, fmt.Errorf("storage: create journal: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Journal{}, fmt.Errorf("storage: create journal: %w", err)
	}

	journal, err := s.GetJournalByID(id)
	if err != nil {
		return Journal{}, err
	}

	return journal, nil
}

func (s *SQLiteStorage) GetJournalByID(id int64) (Journal, error) {
	var j Journal

	err := s.db.QueryRow(
		`SELECT id, name, current_mileage
		FROM journals
		WHERE id = ?
		`, id,
	).Scan(&j.ID, &j.Name, &j.CurrentMileage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Journal{}, errNotFound("journal", id)
		}
		return Journal{}, fmt.Errorf("storage: get journal: %w", err)
	}

	return j, nil
}

func (s *SQLiteStorage) GetAllJournals() ([]Journal, error) {
	rows, err := s.db.Query(
		`SELECT id, name, current_mileage
		FROM journals`,
	)
	if err != nil {
		return nil, fmt.Errorf("storage: get all journals: %w", err)
	}
	defer rows.Close()

	journals := make([]Journal, 0)

	for rows.Next() {
		var j Journal
		if err := rows.Scan(&j.ID, &j.Name, &j.CurrentMileage); err != nil {
			return nil, fmt.Errorf("storage: get all journals: %w", err)
		}
		journals = append(journals, j)
	}

	return journals, nil
}

func (s *SQLiteStorage) UpdateJournalName(id int64, name string) (Journal, error) {
	result, err := s.db.Exec("UPDATE journals SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return Journal{}, fmt.Errorf("storage: update journal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return Journal{}, fmt.Errorf("storage: update journal: %w", err)
	}
	if rows == 0 {
		return Journal{}, errNotFound("journal", id)
	}

	journal, err := s.GetJournalByID(id)
	if err != nil {
		return Journal{}, err
	}

	return journal, nil
}

func (s *SQLiteStorage) UpdateJournalMileage(id int64, mileage int) (Journal, error) {
	result, err := s.db.Exec("UPDATE journals SET current_mileage = ? WHERE id = ?", mileage, id)
	if err != nil {
		return Journal{}, fmt.Errorf("storage: update journal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return Journal{}, fmt.Errorf("storage: update journal: %w", err)
	}
	if rows == 0 {
		return Journal{}, errNotFound("journal", id)
	}

	journal, err := s.GetJournalByID(id)
	if err != nil {
		return Journal{}, err
	}

	return journal, nil
}

func (s *SQLiteStorage) DeleteJournal(id int64) error {
	result, err := s.db.Exec("DELETE FROM journals WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("storage: delete journal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("storage: delete journal: %w", err)
	}

	if rows == 0 {
		return errNotFound("journal", id)
	}

	return nil
}
