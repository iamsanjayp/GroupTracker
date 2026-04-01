package handlers

import (
	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type AttendanceHandler struct {
	attRepo *repository.AttendanceRepo
}

func NewAttendanceHandler(attRepo *repository.AttendanceRepo) *AttendanceHandler {
	return &AttendanceHandler{attRepo: attRepo}
}

// Admins only
func (h *AttendanceHandler) SaveAttendance(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	var req models.SaveAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Date == "" || (req.Session != "morning" && req.Session != "afternoon") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Valid Date and Session (morning/afternoon) required"})
	}

	// Validate hour slots
	for _, rec := range req.Records {
		if req.Session == "morning" && (rec.HourSlot < 1 || rec.HourSlot > 4) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Morning session must only use hours 1-4"})
		}
		if req.Session == "afternoon" && (rec.HourSlot < 5 || rec.HourSlot > 7) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Afternoon session must only use hours 5-7"})
		}
	}

	err := h.attRepo.SaveBulk(teamID, req.Date, req.Records)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save attendance"})
	}

	return c.JSON(fiber.Map{"message": "Attendance saved successfully"})
}

// Admins only
func (h *AttendanceHandler) GetAttendance(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	date := c.Query("date")
	session := c.Query("session")

	if date == "" || (session != "morning" && session != "afternoon") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Valid date and session required"})
	}

	records, err := h.attRepo.GetByDateAndSession(teamID, date, session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch records"})
	}

	// Make sure we return an empty array not null
	if records == nil {
		records = []models.Attendance{}
	}

	return c.JSON(records)
}

// Any member
func (h *AttendanceHandler) LogMissedOTP(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	userID := middleware.GetUserID(c)

	var req models.SaveMissedAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Date == "" || req.HourSlot < 1 || req.HourSlot > 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Valid date and hour slot (1-7) required"})
	}

	err := h.attRepo.LogMissedOTP(teamID, userID, req.Date, req.HourSlot)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to log missed OTP"})
	}

	return c.JSON(fiber.Map{"message": "Missed OTP logged successfully"})
}

// Admins only
func (h *AttendanceHandler) GetMissedOTPExports(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	exports, err := h.attRepo.GetMissedOTPExports(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate exports"})
	}

	if exports == nil {
		exports = []models.MissedAttendanceExport{}
	}

	return c.JSON(exports)
}
