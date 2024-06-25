package models

import (
	"database/sql"
	"time"
)

type Unavailability struct {
	EventId int
	UserId int
	UnavailabilityId int
	AllDay bool
	Start time.Time
	End time.Time
}

type UnavailabilityModel struct {
	DB *sql.DB
}

func (m *UnavailabilityModel) Add(userId, eventId int, start, end string, allDay bool) error {

	stmt := `INSERT INTO unavailabilities (event_id, user_id, all_day, start, end)
			VALUES(?, ?, ?, ?, ?);`


	_, err := m.DB.Exec(stmt, eventId, userId, allDay, start, end)

	if err != nil {
		return err
	}

	return nil
}

func (m *UnavailabilityModel) GetAllUnavailabilities(eventId int) ([]Unavailability, error) {
	
	stmt := `SELECT * FROM unavailabilities WHERE event_id = ?`

	rows, err := m.DB.Query(stmt, eventId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var unavailabilities []Unavailability

	for rows.Next() {
		var u Unavailability

		err := rows.Scan(&u.EventId, &u.UserId, &u.UnavailabilityId, &u.AllDay, &u.Start, &u.End)
		if err != nil {
			return nil, err
		}

		unavailabilities = append(unavailabilities, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return unavailabilities, nil
}