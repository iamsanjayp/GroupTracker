package models

import "time"

type Project struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	TeamID      uint64    `json:"team_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectWithMembers struct {
	Project
	Members []ProjectMember `json:"members"`
	Tasks   []Task          `json:"tasks"`
}

type ProjectMember struct {
	ID              uint64  `json:"id"`
	UserID          uint64  `json:"user_id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	SharePercentage float64 `json:"share_percentage"`
}

type Task struct {
	ID          uint64    `json:"id"`
	ProjectID   uint64    `json:"project_id"`
	TeamID      uint64    `json:"team_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	AssignedTo  *uint64   `json:"assigned_to"`
	AssigneeName *string  `json:"assignee_name"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     *string   `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type AddProjectMemberRequest struct {
	UserID          uint64  `json:"user_id"`
	SharePercentage float64 `json:"share_percentage"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	AssignedTo  *uint64 `json:"assigned_to"`
	Priority    string  `json:"priority"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssignedTo  *uint64 `json:"assigned_to"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"`
}
