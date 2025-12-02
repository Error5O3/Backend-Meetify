package event

import (
	"context"
)

type PublicSlot struct {
	StartTime    string `json:"start_time"`
	ID           int64  `json:"id"`
	NumAvailable int    `json:"num_available"`
}
type PublicGridResponse struct {
	Dates     []string       `json:"dates"`
	TimeSlots [][]PublicSlot `json:"time_slots"`
	Users     []string       `json:"users"`
	NumUsers  int            `json:"numUsers"`
}

type Event struct {
	EventID   int64  `json:"event_id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	CreatedAt string `json:"created_at"`
}

type EventDate struct {
	ID        int64  `json:"id"`
	EventID   int64  `json:"event_id"`
	EventDate string `json:"event_date"`
}

type TimeSlot struct {
	ID          int64  `json:"id"`
	EventDateID int64  `json:"event_date_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type CreateEventRequest struct {
	Name      string   `json:"name" binding:"required"`
	Dates     []string `json:"dates" binding:"required"`
	StartTime string   `json:"start_time" binding:"required"`
	EndTime   string   `json:"end_time" binding:"required"`
}

type CreateEventResponse struct {
	EventID   int64    `json:"event_id"`
	Name      string   `json:"name"`
	Dates     []string `json:"dates"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
}

type EventGridResponse struct {
	EventID   int64       `json:"event_id"`
	Name      string      `json:"name"`
	StartTime string      `json:"start_time"`
	EndTime   string      `json:"end_time"`
	Dates     []DateSlots `json:"dates"`
}

type DateSlots struct {
	Date  string     `json:"date"`
	Slots []SlotInfo `json:"slots"`
}

type SlotInfo struct {
	SlotID         int64  `json:"slot_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	AvailableCount int    `json:"available_count"`
	IsAvailable    bool   `json:"is_available"` // for current user
}

type MarkAvailabilityRequest struct {
	TimeSlotID int64 `json:"time_slot_id" binding:"required"`
}

type VueSlot struct {
	ID             int64    `json:"id"`
	StartTime      string   `json:"start_time"`
	EndTime        string   `json:"end_time"`
	AvailableCount int      `json:"available_count"`
	AvailableUsers []string `json:"available_users"`
}

type Repository interface {
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	CreateEventDate(ctx context.Context, eventDate *EventDate) (*EventDate, error)
	CreateTimeSlot(ctx context.Context, timeSlot *TimeSlot) error
	GetEvent(ctx context.Context, eventID int64) (*Event, error)
	GetEventGrid(ctx context.Context, eventID int64) (*PublicGridResponse, error)
	MarkAvailable(ctx context.Context, userID, timeSlotID int64) error
	UnmarkAvailable(ctx context.Context, userID, timeSlotID int64) error
}

type Service interface {
	CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error)
	GetEvent(ctx context.Context, eventID int64) (*Event, error)
	GetEventGrid(ctx context.Context, eventID int64) (*PublicGridResponse, error)
	MarkAvailable(ctx context.Context, userID int64, req *MarkAvailabilityRequest) error
	UnmarkAvailable(ctx context.Context, userID int64, req *MarkAvailabilityRequest) error
}
