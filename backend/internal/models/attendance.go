package models

import "time"

type Attendance struct {
	ID        uint64    `json:"id"`
	TeamID    uint64    `json:"team_id"`
	UserID    uint64    `json:"user_id"`
	Date      string    `json:"date"` // YYYY-MM-DD
	HourSlot  int       `json:"hour_slot"`
	Status    string    `json:"status"` // 'Present', 'Absent', 'PS Slot', 'Event', 'OnDuty', 'Class'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MissedAttendance struct {
	ID        uint64    `json:"id"`
	TeamID    uint64    `json:"team_id"`
	UserID    uint64    `json:"user_id"`
	Date      string    `json:"date"` // YYYY-MM-DD
	HourSlot  int       `json:"hour_slot"`
	CreatedAt time.Time `json:"created_at"`
}

type SaveAttendanceRequest struct {
	Date       string       `json:"date"`
	Session    string       `json:"session"` // "morning" or "afternoon"
	Records    []Attendance `json:"records"`
}

type MissedAttendanceExport struct {
	Date     string `json:"Date"`
	RollNo   string `json:"Roll No"`
	Name     string `json:"Name"`
	Email    string `json:"Mail Id"`
	HourSlot int    `json:"Hour"`
}

type SaveMissedAttendanceRequest struct {
	Date     string `json:"date"`
	HourSlot int    `json:"hour_slot"`
}
