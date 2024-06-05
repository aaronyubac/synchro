package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
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

func (m *EventModel) GetEvent(id int) (Event, error) {

	stmt := `SELECT * FROM events WHERE event_id = ?`
/* 
	stmt := `SELECT e.event_id, e.event_name, e.event_details
	 FROM events e JOIN users_events ue
	 WHERE event_id = ?` */
	//change stmt so that it checks if user id is linked

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







func (m *EventModel) Join(userID, eventID int) error {

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