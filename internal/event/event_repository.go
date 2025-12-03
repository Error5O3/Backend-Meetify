package event

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateEvent(ctx context.Context, event *Event) (*Event, error) {
	var eventID int64
	query := `INSERT INTO events(name, start_time, end_time) VALUES ($1, $2, $3) RETURNING event_id`
	err := r.db.QueryRowContext(ctx, query, event.Name, event.StartTime, event.EndTime).Scan(&eventID)
	if err != nil {
		return nil, err
	}
	event.EventID = eventID
	return event, nil
}

func (r *repository) CreateEventDate(ctx context.Context, eventDate *EventDate) (*EventDate, error) {
	var id int64
	query := `INSERT INTO event_dates(event_id, event_date) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, eventDate.EventID, eventDate.EventDate).Scan(&id)
	if err != nil {
		return nil, err
	}
	eventDate.ID = id
	return eventDate, nil

}

func (r *repository) CreateTimeSlot(ctx context.Context, timeSlot *TimeSlot) error {
	query := `INSERT INTO time_slots(event_date_id, start_time, end_time) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, timeSlot.EventDateID, timeSlot.StartTime, timeSlot.EndTime)
	return err
}

func (r *repository) GetEvent(ctx context.Context, eventID int64) (*Event, error) {
	query := `SELECT event_id, name, start_time, end_time, created_at 
	          FROM events WHERE event_id = $1`

	var event Event
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.EventID, &event.Name, &event.StartTime, &event.EndTime, &event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// this function will return the availabity of the current user to certain event that has been created before, including:
// the time slot availaiblity with the current users selection

func (r *repository) GetEventGrid(ctx context.Context, eventID int64) (*PublicGridResponse, error) {
	query := `SELECT
				DATE(ed.event_date) AS event_date,
				ts.id AS slot_id,
				ts.start_time,
				COUNT(ua.user_id) AS num_available
			FROM event_dates ed
			JOIN time_slots ts ON ts.event_date_id = ed.id
			LEFT JOIN user_availability ua ON ua.time_slot_id = ts.id
			WHERE ed.event_id = $1
			GROUP BY ed.event_date, ts.id, ts.start_time
			ORDER BY ed.event_date, ts.start_time
				`

	rows, err := r.db.QueryContext(ctx, query, eventID)

	if err != nil {
		return nil, err
	}

	nameQuery := `SELECT name FROM events WHERE event_id = $1`

	var name string
	err1 := r.db.QueryRowContext(ctx, nameQuery, eventID).Scan(&name)

	if err1 != nil {
		return nil, err
	}

	defer rows.Close()
	dateMap := make(map[string][]PublicSlot)
	var dates []string
	dateOrder := make(map[string]int)

	for rows.Next() {
		var date, startTime string
		var slotID int64
		var numAvailable int
		err := rows.Scan(&date, &slotID, &startTime, &numAvailable)
		if err != nil {
			return nil, err
		}
		if _, exists := dateMap[date]; !exists {
			dateOrder[date] = len(dates)
			dates = append(dates, date)
			dateMap[date] = []PublicSlot{}
		}
		dateMap[date] = append(dateMap[date], PublicSlot{
			StartTime:    startTime,
			ID:           slotID,
			NumAvailable: numAvailable,
		})
	}

	timeSlots := make([][]PublicSlot, len(dates))
	for date, slots := range dateMap {
		index := dateOrder[date]
		timeSlots[index] = slots
	}

	usersQuery := `SELECT DISTINCT u.username
							FROM users u
							JOIN user_availability ua ON ua.user_id = u.id
							JOIN time_slots ts ON ts.id = ua.time_slot_id
							JOIN event_dates ed ON ed.id = ts.event_date_id
							WHERE ed.event_id = $1
							ORDER BY u.username`
	userRows, err := r.db.QueryContext(ctx, usersQuery, eventID)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var users []string
	for userRows.Next() {
		var username string
		userRows.Scan(&username)
		users = append(users, username)
	}

	return &PublicGridResponse{
		Dates:     dates,
		TimeSlots: timeSlots,
		Users:     users,
		NumUsers:  len(users),
		EventName: name,
	}, nil
}

// basically for this query we need to do a for loop because we are given a set of id slots we wanna add for an specific user	return &EventGridResponse{
// 		EventID:   event.EventID,
// 		Name:      event.Name,
// 		StartTime: event.StartTime,
// 		EndTime:   event.EndTime,
// 		Dates:     dates,
// 	}, nil
// }

func (r *repository) MarkAvailable(ctx context.Context, userID, timeSlotID int64) error {

	// this will be triggered every time you dealocate the time slot of a userquery := `
	query := `INSERT INTO user_availability (user_id, time_slot_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, timeSlotID)
	return err
}

func (r *repository) UnmarkAvailable(ctx context.Context, userID, timeSlotID int64) error {
	query := `DELETE FROM user_availability WHERE user_id = $1 AND time_slot_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, timeSlotID)
	return err
}
