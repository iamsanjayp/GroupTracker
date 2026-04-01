package handlers

import (
	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type PointsHandler struct {
	pointsRepo *repository.PointsRepo
}

func NewPointsHandler(pointsRepo *repository.PointsRepo) *PointsHandler {
	return &PointsHandler{pointsRepo: pointsRepo}
}

func (h *PointsHandler) GetMyPoints(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)

	points, err := h.pointsRepo.GetUserPoints(userID, teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get points"})
	}

	// Get PS records too
	psRecords, _ := h.pointsRepo.GetPSRecords(userID, teamID)

	return c.JSON(fiber.Map{
		"points":     points,
		"ps_records": psRecords,
	})
}

func (h *PointsHandler) GetTeamLeaderboard(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	leaderboard, err := h.pointsRepo.GetTeamLeaderboard(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get leaderboard"})
	}

	return c.JSON(leaderboard)
}

func (h *PointsHandler) AddPSRecord(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)

	var req models.CreatePSRecordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.CourseName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reason/Course name is required"})
	}

	if req.RewardPoints == 0 && req.ActivityPoints == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one point value must be non-zero"})
	}

	if err := h.pointsRepo.CreatePSRecord(userID, teamID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add PS record"})
	}

	// Recalculate
	h.pointsRepo.RecalculateUser(userID, teamID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "PS record added"})
}

func (h *PointsHandler) GetPointsHistory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	history, err := h.pointsRepo.GetPointsHistory(userID, teamID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get point history"})
	}

	return c.JSON(history)
}

// Admins only
func (h *PointsHandler) AddBulkPoints(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	var req models.BulkPointUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	if len(req.Records) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No records to process"})
	}

	resp, err := h.pointsRepo.BulkAddPoints(teamID, req.Records)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to bulk upload points"})
	}

	return c.JSON(resp)
}
