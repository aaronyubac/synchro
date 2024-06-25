package models

import (
	"database/sql"
	"time"
)

type Unavailability struct {
	EventId int
	UserId int
	UnavailabilityId int
	Date time.Time
	AllDay bool
	Start time.Time
	End time.Time
}

type UnavailabilityModel struct {
	DB *sql.DB
}

func (m *UnavailabilityModel) Add(userId, eventId int, date time.Time, start, end string, allDay bool) error {

	stmt := `INSERT INTO unavailabilities (user_id, event_id, date, all_day, start, end)
	VALUES (?, ?, ?, ?, ?, ?)`

	var err error

	if allDay {
		_, err = m.DB.Exec(stmt, userId, eventId, date, allDay, nil, nil)
	} else {
		_, err = m.DB.Exec(stmt, userId, eventId, date, allDay, start, end)
	}

	if err != nil {
		return err
	}

	return nil
}