package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/config"
	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type TeamHandler struct {
	cfg      *config.Config
	teamRepo *repository.TeamRepo
	userRepo *repository.UserRepo
}

func NewTeamHandler(cfg *config.Config, teamRepo *repository.TeamRepo, userRepo *repository.UserRepo) *TeamHandler {
	return &TeamHandler{cfg: cfg, teamRepo: teamRepo, userRepo: userRepo}
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Team name is required"})
	}

	// Check if user already in a team
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find user"})
	}
	if user.TeamID != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You already belong to a team"})
	}

	team, err := h.teamRepo.Create(req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create team"})
	}

	// Make creator the captain and immediately approved
	if err := h.userRepo.UpdateTeamAndRole(userID, team.ID, "captain", "approved"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to assign captain role"})
	}

	// Generate new tokens with updated team_id and role
	accessToken, err := middleware.GenerateAccessToken(h.cfg, userID, team.ID, "captain")
	if err != nil {
		fmt.Printf("Warning: failed to generate new token: %v\n", err)
	}

	refreshToken, expiry, err := middleware.GenerateRefreshToken(h.cfg, userID)
	if err != nil {
		fmt.Printf("Warning: failed to generate refresh token: %v\n", err)
	}
	h.userRepo.SaveRefreshToken(userID, refreshToken, expiry)

	// Reload user
	user, _ = h.userRepo.FindByID(userID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"team":          team,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user.ToResponse(),
	})
}

func (h *TeamHandler) JoinTeam(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req struct {
		InviteCode string `json:"invite_code"`
	}
	if err := c.BodyParser(&req); err != nil || req.InviteCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invite code is required"})
	}

	// Check if user already in a team
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find user"})
	}
	if user.TeamID != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You already belong to a team"})
	}

	team, err := h.teamRepo.FindByInviteCode(req.InviteCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid invite code"})
	}

	// Join as pending member
	if err := h.userRepo.UpdateTeamAndRole(userID, team.ID, "member", "pending"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to join team"})
	}

	// Generate new tokens with updated team_id
	accessToken, _ := middleware.GenerateAccessToken(h.cfg, userID, team.ID, "member")
	refreshToken, expiry, _ := middleware.GenerateRefreshToken(h.cfg, userID)
	h.userRepo.SaveRefreshToken(userID, refreshToken, expiry)

	user, _ = h.userRepo.FindByID(userID)

	return c.JSON(fiber.Map{
		"message":       "Successfully joined team",
		"team":          team,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user.ToResponse(),
	})
}

func (h *TeamHandler) GetMyTeam(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	team, err := h.teamRepo.FindByID(teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Team not found"})
	}

	memberCount, _ := h.teamRepo.GetMemberCount(teamID)

	return c.JSON(fiber.Map{
		"team":         team,
		"member_count": memberCount,
	})
}

func (h *TeamHandler) GetMembers(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	members, err := h.userRepo.GetTeamMembers(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get members"})
	}

	return c.JSON(members)
}

func (h *TeamHandler) UpdateMemberRole(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	memberID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Role is required"})
	}

	validRoles := map[string]bool{
		"captain": true, "vice_captain": true, "manager": true, "strategist": true, "member": true,
	}
	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
	}

	if err := h.userRepo.UpdateRole(uint64(memberID), teamID, req.Role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update role"})
	}

	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

func (h *TeamHandler) RemoveMember(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	userID := middleware.GetUserID(c)
	memberID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	if uint64(memberID) == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You cannot remove yourself"})
	}

	if err := h.userRepo.RemoveFromTeam(uint64(memberID), teamID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove member"})
	}

	return c.JSON(fiber.Map{"message": "Member removed successfully"})
}

func (h *TeamHandler) GetPendingMembers(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)

	members, err := h.userRepo.GetTeamMembers(teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get members"})
	}

	// Filter for pending members
	var pending []models.TeamMember
	for _, m := range members {
		if m.JoinStatus != nil && *m.JoinStatus == "pending" {
			pending = append(pending, m)
		}
	}

	return c.JSON(pending)
}

func (h *TeamHandler) ApproveMember(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	memberID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	// Update join_status to 'approved'. We can reuse UpdateTeamAndRole or create a new method.
	// Since we know they are a 'member' currently, it's safer to just fetch them or use a dedicated query.
	// We'll update just their join_status.
	// Actually, let's just create an UpdateJoinStatus method on userRepo or execute here.
	user, err := h.userRepo.FindByID(uint64(memberID))
	if err != nil || user.TeamID == nil || *user.TeamID != teamID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Member not found in team"})
	}

	err = h.userRepo.UpdateTeamAndRole(uint64(memberID), teamID, user.Role, "approved")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to approve member"})
	}

	return c.JSON(fiber.Map{"message": "Member approved successfully"})
}
