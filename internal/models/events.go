package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type Event struct {
	ID string
	Name string
	Details string
}

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) Create(name, details string) (string, error) {

	code, err := generateHexCode(6)
	if err != nil {
		return "", err
	}

	stmt := `INSERT INTO events (event_id, event_name, event_details) VALUES(?, ?, ?)`

	_, err = m.DB.Exec(stmt, code, name, details)
	if err != nil {
		return "", err
	}

	return code, err
}

func (m *EventModel) GetEvent(userID int, eventID string) (Event, error) {

	stmt := `SELECT e.event_id, e.event_name, e.event_details
	 FROM events e JOIN users_events ue
	 ON e.event_id = ue.event_id
	 WHERE ue.user_id = ? AND e.event_id = ?`

	row := m.DB.QueryRow(stmt, userID, eventID)

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

func (m *EventModel) GetUserEvents(userID int) ([]Event, error) {

	stmt := `SELECT e.event_id, e.event_name, e.event_details 
			 FROM events e JOIN users_events ue
			 ON e.event_id = ue.event_id
			 WHERE ue.user_id = ?`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var e Event

		err := rows.Scan(&e.ID, &e.Name, &e.Details)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil

}







func (m *EventModel) Join(userID int, eventID string) error {

	stmt := `INSERT INTO users_events (user_id, event_id) VALUES (?, ?)`

	_, err := m.DB.Exec(stmt, userID, eventID)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1452 {
				return ErrNoRecord
			}
			if mySQLError.Number == 1062 {
				return ErrDuplicateEvent
			}
		}
	}

	return nil
}


func generateHexCode(n int) (string, error) {
	bytes := make([]byte, (n+1) / 2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}