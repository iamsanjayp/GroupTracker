package models

import "time"

type Skill struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	TeamID   uint64 `json:"team_id"`
}

type UserSkill struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	SkillID   uint64    `json:"skill_id"`
	SkillName string    `json:"skill_name"`
	Category  string    `json:"category"`
	TeamID    uint64    `json:"team_id"`
	Level     int       `json:"level"`
	Validated bool      `json:"validated"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
