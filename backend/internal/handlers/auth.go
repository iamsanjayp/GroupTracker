package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"grouptracker/internal/config"
	"grouptracker/internal/middleware"
	"grouptracker/internal/models"
	"grouptracker/internal/repository"
)

type AuthHandler struct {
	cfg      *config.Config
	userRepo *repository.UserRepo
}

func NewAuthHandler(cfg *config.Config, userRepo *repository.UserRepo) *AuthHandler {
	return &AuthHandler{cfg: cfg, userRepo: userRepo}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" || req.Name == "" || req.RollNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email, password, name, and Roll No are required"})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
	}

	// Check if email exists
	existing, _ := h.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	approvedStatus := "approved"
	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		Role:         "member",
		RollNo:       &req.RollNo,
		JoinStatus:   &approvedStatus,
	}

	userID, err := h.userRepo.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	user.ID = userID
	return h.generateAuthResponse(c, user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if user.PasswordHash == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "This account does not have a password set. Please use the 'Forgot Password' flow or contact support."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return h.generateAuthResponse(c, user)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req models.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Parse the refresh token to get user_id
	claims, err := middleware.ParseRefreshToken(h.cfg, req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	// Validate refresh token exists in DB
	valid, err := h.userRepo.ValidateRefreshToken(userID, req.RefreshToken)
	if err != nil || !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token expired or revoked"})
	}

	// Get fresh user data
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Delete old token
	h.userRepo.DeleteRefreshToken(userID, req.RefreshToken)

	return h.generateAuthResponse(c, user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	h.userRepo.DeleteAllRefreshTokens(userID)
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user.ToResponse())
}

func (h *AuthHandler) generateAuthResponse(c *fiber.Ctx, user *models.User) error {
	var teamID uint64
	if user.TeamID != nil {
		teamID = *user.TeamID
	}

	accessToken, err := middleware.GenerateAccessToken(h.cfg, user.ID, teamID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	refreshToken, expiry, err := middleware.GenerateRefreshToken(h.cfg, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	// Save refresh token hash in DB
	if err := h.userRepo.SaveRefreshToken(user.ID, refreshToken, expiry); err != nil {
		fmt.Printf("Warning: failed to save refresh token: %v\n", err)
	}

	return c.JSON(models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	})
}
