package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type ActivityHandler struct {
	activityRepo *repository.ActivityRepo
	pointsRepo   *repository.PointsRepo
}

func NewActivityHandler(activityRepo *repository.ActivityRepo, pointsRepo *repository.PointsRepo) *ActivityHandler {
	return &ActivityHandler{activityRepo: activityRepo, pointsRepo: pointsRepo}
}

func (h *ActivityHandler) GetDayLog(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)
	date := c.Query("date", time.Now().Format("2006-01-02"))

	activities, err := h.activityRepo.GetDayLog(userID, teamID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get activities"})
	}

	actTotal, rewTotal, _ := h.activityRepo.GetUserTodayPoints(userID, teamID, date)

	return c.JSON(models.DayLog{
		Date:                date,
		Activities:          activities,
		TotalActivityPoints: actTotal,
		TotalRewardPoints:   rewTotal,
	})
}

func (h *ActivityHandler) BulkSave(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)

	var req models.BulkActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	for i := range req.Activities {
		entry := &req.Activities[i]

		// Validate hour slot
		if entry.HourSlot < 1 || entry.HourSlot > 7 {
			continue
		}

		// Points are set by the frontend (auto-suggested or manual)
		// 0 points is valid (some activities don't carry points)

		if err := h.activityRepo.UpsertActivity(userID, teamID, req.Date, *entry); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save activity for hour " + string(rune('0'+entry.HourSlot)),
			})
		}
	}

	// Recalculate points
	h.pointsRepo.RecalculateUser(userID, teamID)

	return c.JSON(fiber.Map{"message": "Activities saved successfully"})
}

func (h *ActivityHandler) GetTeamActivities(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	date := c.Query("date", time.Now().Format("2006-01-02"))

	activities, err := h.activityRepo.GetTeamActivities(teamID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get team activities"})
	}

	return c.JSON(fiber.Map{
		"date":       date,
		"activities": activities,
	})
}
