package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type ProjectHandler struct {
	projectRepo *repository.ProjectRepo
}

func NewProjectHandler(projectRepo *repository.ProjectRepo) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo}
}

func (h *ProjectHandler) List(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	projects, err := h.projectRepo.GetByTeam(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get projects"})
	}

	// Enrich with member count
	type ProjectWithCount struct {
		models.Project
		MemberCount int `json:"member_count"`
		TaskCount   int `json:"task_count"`
	}

	var enriched []ProjectWithCount
	for _, p := range projects {
		members, _ := h.projectRepo.GetMembers(p.ID, teamID)
		tasks, _ := h.projectRepo.GetTasks(p.ID, teamID)
		enriched = append(enriched, ProjectWithCount{
			Project:     p,
			MemberCount: len(members),
			TaskCount:   len(tasks),
		})
	}

	return c.JSON(enriched)
}

func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	projectID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	project, err := h.projectRepo.GetByID(projectID, teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	members, _ := h.projectRepo.GetMembers(projectID, teamID)
	tasks, _ := h.projectRepo.GetTasks(projectID, teamID)

	return c.JSON(models.ProjectWithMembers{
		Project: *project,
		Members: members,
		Tasks:   tasks,
	})
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	var req models.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project name is required"})
	}

	// Enforce 3-5 project limit
	count, _ := h.projectRepo.GetProjectCount(teamID)
	if count >= 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Maximum 5 projects per team"})
	}

	id, err := h.projectRepo.Create(req.Name, req.Description, teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id, "message": "Project created"})
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	projectID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req models.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.projectRepo.Update(projectID, teamID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update project"})
	}

	return c.JSON(fiber.Map{"message": "Project updated"})
}

func (h *ProjectHandler) AddMember(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	projectID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req models.AddProjectMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.projectRepo.AddMember(projectID, req.UserID, teamID, req.SharePercentage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add member"})
	}

	return c.JSON(fiber.Map{"message": "Member added to project"})
}

func (h *ProjectHandler) RemoveMember(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	projectID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}
	userID, err := strconv.ParseUint(c.Params("uid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.projectRepo.RemoveMember(projectID, userID, teamID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove member"})
	}

	return c.JSON(fiber.Map{"message": "Member removed from project"})
}

func (h *ProjectHandler) CreateTask(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	projectID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Task title is required"})
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	id, err := h.projectRepo.CreateTask(projectID, teamID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id, "message": "Task created"})
}

func (h *ProjectHandler) UpdateTask(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	taskID, err := strconv.ParseUint(c.Params("tid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.projectRepo.UpdateTask(taskID, teamID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.JSON(fiber.Map{"message": "Task updated"})
}
