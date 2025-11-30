package event

import(
	"context"
	"time"
)

type service struct{
	Repository
}

// param elements that are going inside the service struct
func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateEvent(c context.Context, req *CreateEventRequest, userID int64) (*CreateEventResponse, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	event := &Event{
			Name : req.Name,
			StartTime : req.StartTime,
			EndTime : req.EndTime,
			CreatedBy : userID,
	}
	createdEvent, err := s.Repository.CreateEvent(ctx, event)
	if err != nil{
		return nil, err
	}
	for _, dateStr := range req.Dates{
		eventDate := &EventDate{
			EventID: createdEvent.EventID,
			EventDate: dateStr,
		}
		createdDate, err := s.Repository.CreateEventDate(ctx, eventDate)
		if err != nil{
			return nil, err
		}
		if err := s.generateTimeSlots(ctx, createdDate.ID, req.StartTime, req.EndTime); err != nil{
			return nil, err
		}

	}
	res := &CreateEventResponse{
		EventID: createdEvent.EventID,
		Name: createdEvent.Name,
		Dates: req.Dates,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}
	
	return res, nil
}

func (s *service) generateTimeSlots(ctx context.Context, eventDateID int64, startTime, endTime string) error {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return err
	}
	
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return err
	}

	current := start
	for current.Before(end) {
		slotStart := current.Format("15:04")
		current = current.Add(1 * time.Hour)
		slotEnd := current.Format("15:04")

		timeSlot := &TimeSlot{
			EventDateID: eventDateID,
			StartTime:   slotStart,
			EndTime:     slotEnd,
		}

		if err := s.Repository.CreateTimeSlot(ctx, timeSlot); err != nil {
			return err
		}
	}
	
	return nil
}


func (s *service) GetEvent(c context.Context, eventID int64) (*Event, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	return s.Repository.GetEvent(ctx, eventID)
}

func (s *service) GetEventGrid(c context.Context, eventID int64, userID int64) (*EventGridResponse, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	return s.Repository.GetEventGrid(ctx, eventID, userID)
}

func (s *service) MarkAvailable(c context.Context, userID int64, req *MarkAvailabilityRequest) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	return s.Repository.MarkAvailable(ctx, userID, req.TimeSlotID)
}

func (s *service) UnmarkAvailable(c context.Context, userID int64, req *MarkAvailabilityRequest) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	return s.Repository.UnmarkAvailable(ctx, userID, req.TimeSlotID)
}