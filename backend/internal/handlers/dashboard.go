package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"grouptracker/internal/middleware"
	"grouptracker/internal/repository"
)

type DashboardHandler struct {
	activityRepo *repository.ActivityRepo
	projectRepo  *repository.ProjectRepo
	pointsRepo   *repository.PointsRepo
	userRepo     *repository.UserRepo
	attRepo      *repository.AttendanceRepo
}

func NewDashboardHandler(
	activityRepo *repository.ActivityRepo,
	projectRepo *repository.ProjectRepo,
	pointsRepo *repository.PointsRepo,
	userRepo *repository.UserRepo,
	attRepo *repository.AttendanceRepo,
) *DashboardHandler {
	return &DashboardHandler{
		activityRepo: activityRepo,
		projectRepo:  projectRepo,
		pointsRepo:   pointsRepo,
		userRepo:     userRepo,
		attRepo:      attRepo,
	}
}

func (h *DashboardHandler) MemberDashboard(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	teamID := middleware.GetTeamID(c)
	today := time.Now().Format("2006-01-02")
	month := time.Now().Format("2006-01")

	// Today's points
	todayAct, todayRew, _ := h.activityRepo.GetUserTodayPoints(userID, teamID, today)

	// Today's log
	activities, _ := h.activityRepo.GetDayLog(userID, teamID, today)
	hoursLogged := len(activities)

	// Monthly stats
	loggedDays, _ := h.activityRepo.GetUserLoggedDays(userID, teamID, month)

	// Total points
	points, _ := h.pointsRepo.GetUserPoints(userID, teamID)

	// Attendance stats
	attendancePct, _ := h.attRepo.GetUserAttendanceStats(teamID, userID)

	return c.JSON(fiber.Map{
		"today": fiber.Map{
			"date":            today,
			"hours_logged":    hoursLogged,
			"activity_points": todayAct,
			"reward_points":   todayRew,
		},
		"month": fiber.Map{
			"logged_days": loggedDays,
		},
		"total": fiber.Map{
			"activity_points": points.TotalActivity,
			"reward_points":   points.TotalReward,
			"total_points":    points.TotalActivity + points.TotalReward,
		},
		"attendance_percentage": attendancePct,
	})
}

func (h *DashboardHandler) AdminDashboard(c *fiber.Ctx) error {
	teamID := middleware.GetTeamID(c)
	today := time.Now().Format("2006-01-02")

	// Team activity today
	activeToday, totalMembers, _ := h.activityRepo.GetTeamActiveToday(teamID, today)

	// Project stats
	projectCount, _ := h.projectRepo.GetProjectCount(teamID)
	totalTasks, doneTasks, _ := h.projectRepo.GetTaskStats(teamID)

	// Leaderboard
	leaderboard, _ := h.pointsRepo.GetTeamLeaderboard(teamID)

	// Members
	members, _ := h.userRepo.GetTeamMembers(teamID)

	return c.JSON(fiber.Map{
		"today": fiber.Map{
			"date":          today,
			"active_today":  activeToday,
			"total_members": totalMembers,
		},
		"projects": fiber.Map{
			"total":           projectCount,
			"total_tasks":     totalTasks,
			"completed_tasks": doneTasks,
		},
		"leaderboard": leaderboard,
		"members":     members,
	})
}
