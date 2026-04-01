package models

import "time"

type Activity struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"`
	TeamID         uint64    `json:"team_id"`
	ActivityDate   string    `json:"activity_date"`
	HourSlot       int       `json:"hour_slot"`
	ActivityType   string    `json:"activity_type"`
	Description    string    `json:"description"`
	ActivityPoints float64   `json:"activity_points"`
	RewardPoints   float64   `json:"reward_points"`
	ProjectID      *uint64   `json:"project_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ActivityEntry struct {
	HourSlot       int     `json:"hour_slot"`
	ActivityType   string  `json:"activity_type"`
	Description    string  `json:"description"`
	ActivityPoints float64 `json:"activity_points"`
	RewardPoints   float64 `json:"reward_points"`
	ProjectID      *uint64 `json:"project_id"`
}

type BulkActivityRequest struct {
	Date       string          `json:"date"`
	Activities []ActivityEntry `json:"activities"`
}

type DayLog struct {
	Date       string     `json:"date"`
	Activities []Activity `json:"activities"`
	TotalActivityPoints float64 `json:"total_activity_points"`
	TotalRewardPoints   float64 `json:"total_reward_points"`
}

// Default points per activity type when auto-calculating
var DefaultActivityPoints = map[string]float64{
	"project_work":        1.0,
	"ps_slot":             1.0,
	"self_study":          0.75,
	"event":               1.5,
	"class_participation": 0.5,
}
