package models

import "time"

type Team struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	InviteCode string    `json:"invite_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type JoinTeamRequest struct {
	InviteCode string `json:"invite_code"`
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

type TeamMember struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	AvatarURL  *string `json:"avatar_url"`
	IsActive   bool    `json:"is_active"`
	RollNo     *string `json:"roll_no"`
	JoinStatus *string `json:"join_status"`
}
