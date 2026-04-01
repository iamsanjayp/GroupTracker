package models

import "time"

type User struct {
	ID           uint64    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url"`
	TeamID       *uint64   `json:"team_id"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	RollNo       *string   `json:"roll_no"`
	JoinStatus   *string   `json:"join_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID         uint64  `json:"id"`
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	AvatarURL  *string `json:"avatar_url"`
	TeamID     *uint64 `json:"team_id"`
	Role       string  `json:"role"`
	IsActive   bool    `json:"is_active"`
	RollNo     *string `json:"roll_no"`
	JoinStatus string  `json:"join_status"`
}

func (u *User) GetJoinStatus() string {
	if u.JoinStatus != nil {
		return *u.JoinStatus
	}
	return "approved"
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:  u.AvatarURL,
		TeamID:     u.TeamID,
		Role:       u.Role,
		IsActive:   u.IsActive,
		RollNo:     u.RollNo,
		JoinStatus: u.GetJoinStatus(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	RollNo   string `json:"roll_no"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
