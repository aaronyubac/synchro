package models

import (
	"database/sql"
	"errors"
)

type Event struct {
	ID int
	Name string
	Details string
	// Unavailabilities []time.Time
}

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) Create(name, details string) (int, error) {

	stmt := `INSERT INTO events (event_name, event_details) VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, name, details)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *EventModel) Get(id int) (Event, error) {

	stmt := `SELECT * FROM events WHERE event_id = ?`

	row := m.DB.QueryRow(stmt, id)

	var e Event

	err := row.Scan(&e.ID, &e.Name, &e.Details)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Event{}, ErrNoRecord
		} else {
			return Event{}, err
		}
	}

	return e, nil

}

/* func (m *EventModel) Join(id int) error {

	// stmt := `INSERT INTO users_events (user_id, event_id) VALUES (?, ?)`

} */