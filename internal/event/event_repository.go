package event

import(
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct{
	db DBTX
}

func NewRepository(db DBTX) *repository{
	return &repository{
		db : db,
	} 
}

func(r *repository) CreateEvent(ctx context.Context, event *Event) (*Event, error){
	var eventID int64
	query := `INSERT INTO events(name, start_time, end_time,created_by) VALUES ($1, $2, $3, $4) RETURNING event_id`
	err := r.db.QueryRowContext(ctx, query, event.Name, event.StartTime, event.EndTime,event.CreatedBy).Scan(&eventID)
	if err != nil{
		return nil, err
	}
	event.EventID = eventID
	return  event, nil
}

func(r *repository) CreateEventDate(ctx context.Context, eventDate *EventDate) (*EventDate, error){
	var id int64
	query := `INSERT INTO event_dates(event_id, event_date) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, eventDate.EventID, eventDate.EventDate).Scan(&id)
	if err != nil{
		return nil, err
	}
	eventDate.ID = id
	return eventDate, nil
	
}

func(r *repository) CreateTimeSlot(ctx context.Context, timeSlot *TimeSlot) error {
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

// func (r *repository) GetEventGrid(ctx context.Context, eventID int64, userID int64) (*EventGridResponse, error) {
// 	// Get event details
// 	event, err := r.GetEvent(ctx, eventID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get all slots with availability counts and user's selections
// 	query := `
// 		SELECT 
// 			ed.event_date,
// 			ts.id as slot_id,
// 			ts.start_time, 
// 			ts.end_time,
// 			COUNT(DISTINCT ua.user_id) as available_count,
// 			CASE WHEN ua_user.user_id IS NOT NULL THEN true ELSE false END as is_available
// 		FROM event_dates ed
// 		JOIN time_slots ts ON ts.event_date_id = ed.id
// 		LEFT JOIN user_availability ua ON ua.time_slot_id = ts.id
// 		LEFT JOIN user_availability ua_user ON ua_user.time_slot_id = ts.id AND ua_user.user_id = $2
// 		WHERE ed.event_id = $1
// 		GROUP BY ed.event_date, ts.id, ts.start_time, ts.end_time, ua_user.user_id
// 		ORDER BY ed.event_date, ts.start_time
// 	`
	
// 	rows, err := r.db.QueryContext(ctx, query, eventID, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Build the response
// 	dateMap := make(map[string]*DateSlots)
// 	var dates []DateSlots

// 	for rows.Next() {
// 		var date string
// 		var slot SlotInfo
		
// 		err := rows.Scan(&date, &slot.SlotID, &slot.StartTime, &slot.EndTime, 
// 			&slot.AvailableCount, &slot.IsAvailable)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if _, exists := dateMap[date]; !exists {
// 			dateMap[date] = &DateSlots{
// 				Date:  date,
// 				Slots: []SlotInfo{},
// 			}
// 			dates = append(dates, *dateMap[date])
// 		}
// 		dateMap[date].Slots = append(dateMap[date].Slots, slot)
// 	}

// 	// Update dates slice with complete slot info
// 	for i := range dates {
// 		if ds, ok := dateMap[dates[i].Date]; ok {
// 			dates[i].Slots = ds.Slots
// 		}
// 	}

// 	return &EventGridResponse{
// 		EventID:   event.EventID,
// 		Name:      event.Name,
// 		StartTime: event.StartTime,
// 		EndTime:   event.EndTime,
// 		Dates:     dates,
// 	}, nil
// }

func (r *repository) MarkAvailable(ctx context.Context, userID, timeSlotID int64) error {
	query := `
			INSERT INTO user_availability (user_id, time_slot_id)
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