package handlers

import (
	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/repository"
)

type SkillHandler struct {
	skillRepo *repository.SkillRepo
}

func NewSkillHandler(skillRepo *repository.SkillRepo) *SkillHandler {
	return &SkillHandler{skillRepo: skillRepo}
}

// GetFixedSkills returns the hardcoded list of available skills per category
func (h *SkillHandler) GetFixedSkills(c *fiber.Ctx) error {
	return c.JSON(repository.FixedSkills)
}

// GetMySkills returns the current user's selected skills
func (h *SkillHandler) GetMySkills(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)

	skills, err := h.skillRepo.GetUserSkills(userID, teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get skills"})
	}

	return c.JSON(fiber.Map{"skills": skills})
}

// SetMySkills sets the current user's skills (first time only, unless captain)
func (h *SkillHandler) SetMySkills(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)
	role := middleware.GetRole(c)

	// Check if already set — only captain can change after first set
	hasSkills, _ := h.skillRepo.HasSkillsSet(userID, teamID)
	if hasSkills && role != "captain" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Skills already set. Only the captain can modify skills.",
		})
	}

	var req repository.SetSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.skillRepo.SetUserSkills(userID, teamID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Skills saved successfully"})
}

// SetMemberSkills lets the captain set/change any team member's skills
func (h *SkillHandler) SetMemberSkills(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	memberID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	var req repository.SetSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.skillRepo.SetUserSkills(uint64(memberID), teamID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Member skills updated"})
}

// GetTeamSkills returns all team members' skills (admin view)
func (h *SkillHandler) GetTeamSkills(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	skills, err := h.skillRepo.GetTeamSkills(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get team skills"})
	}

	return c.JSON(skills)
}
