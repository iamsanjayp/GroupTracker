package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"grouptracker/internal/config"
	"grouptracker/internal/database"
	"grouptracker/internal/handlers"
	"grouptracker/internal/middleware"
	"grouptracker/internal/repository"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to database
	db := database.Connect(cfg)
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	activityRepo := repository.NewActivityRepo(db)
	projectRepo := repository.NewProjectRepo(db)
	pointsRepo := repository.NewPointsRepo(db)
	attRepo := repository.NewAttendanceRepo(db)
	skillRepo := repository.NewSkillRepo(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, userRepo)
	teamHandler := handlers.NewTeamHandler(cfg, teamRepo, userRepo)
	dashboardHandler := handlers.NewDashboardHandler(activityRepo, projectRepo, pointsRepo, userRepo, attRepo)
	activityHandler := handlers.NewActivityHandler(activityRepo, pointsRepo)
	projectHandler := handlers.NewProjectHandler(projectRepo)
	pointsHandler := handlers.NewPointsHandler(pointsRepo)
	attHandler := handlers.NewAttendanceHandler(attRepo)
	skillHandler := handlers.NewSkillHandler(skillRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:        "GroupTracker API",
		ReadBufferSize: 16384,
		ErrorHandler:   customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "grouptracker"})
	})

	// API routes
	api := app.Group("/api")

	// ── Auth (public) ──────────────────────────
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/google", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)
	auth.Post("/google/callback", authHandler.GoogleCallback)
	auth.Post("/refresh", authHandler.Refresh)

	// ── Auth (protected) ───────────────────────
	authProtected := auth.Group("", middleware.AuthMiddleware(cfg))
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Get("/me", authHandler.Me)

	// ── Protected routes (require auth) ────────
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Teams — creating/joining doesn't require team membership
	protected.Post("/teams", teamHandler.CreateTeam)
	protected.Post("/teams/join", teamHandler.JoinTeam)

	// ── Routes requiring team membership ───────
	teamProtected := protected.Group("", middleware.TenantMiddleware())

	// Teams
	teamProtected.Get("/teams/me", teamHandler.GetMyTeam)
	teamProtected.Get("/teams/members", teamHandler.GetMembers)

	// Admin-only team routes
	teamAdmin := teamProtected.Group("", middleware.RequireAdmin())
	teamAdmin.Put("/teams/members/:id/role", teamHandler.UpdateMemberRole)
	teamAdmin.Delete("/teams/members/:id", middleware.RequireRole("captain", "vice_captain"), teamHandler.RemoveMember)

	// Pending members (Captain / VC only)
	teamProtected.Get("/teams/pending", middleware.RequireRole("captain", "vice_captain"), teamHandler.GetPendingMembers)
	teamProtected.Put("/teams/members/:id/approve", middleware.RequireRole("captain", "vice_captain"), teamHandler.ApproveMember)

	// Dashboard
	teamProtected.Get("/dashboard/member", dashboardHandler.MemberDashboard)
	teamProtected.Get("/dashboard/admin", middleware.RequireAdmin(), dashboardHandler.AdminDashboard)

	// Activities
	teamProtected.Get("/activities", activityHandler.GetDayLog)
	teamProtected.Post("/activities/bulk", activityHandler.BulkSave)
	teamProtected.Get("/activities/team", middleware.RequireAdmin(), activityHandler.GetTeamActivities)

	// Projects
	teamProtected.Get("/projects", projectHandler.List)
	teamProtected.Get("/projects/:id", projectHandler.GetByID)

	projectAdmin := teamProtected.Group("/projects", middleware.RequireAdmin())
	projectAdmin.Post("/", projectHandler.Create)
	projectAdmin.Put("/:id", projectHandler.Update)
	projectAdmin.Post("/:id/members", projectHandler.AddMember)
	projectAdmin.Delete("/:id/members/:uid", projectHandler.RemoveMember)
	projectAdmin.Post("/:id/tasks", projectHandler.CreateTask)

	// Task update — any assigned member or admin
	teamProtected.Put("/projects/:id/tasks/:tid", projectHandler.UpdateTask)

	// Points
	teamProtected.Get("/points/me", pointsHandler.GetMyPoints)
	teamProtected.Get("/points/history", pointsHandler.GetPointsHistory)
	teamProtected.Get("/points/team", middleware.RequireAdmin(), pointsHandler.GetTeamLeaderboard)
	teamProtected.Post("/points/ps", pointsHandler.AddPSRecord)
	teamProtected.Post("/points/bulk", middleware.RequireAdmin(), pointsHandler.AddBulkPoints)

	// Attendance
	teamProtected.Post("/attendance", middleware.RequireAdmin(), attHandler.SaveAttendance)
	teamProtected.Get("/attendance", middleware.RequireAdmin(), attHandler.GetAttendance)
	
	// Missed Attendance
	teamProtected.Post("/attendance/missed", attHandler.LogMissedOTP)
	teamProtected.Get("/attendance/missed/exports", middleware.RequireAdmin(), attHandler.GetMissedOTPExports)

	// Skills
	teamProtected.Get("/skills/options", skillHandler.GetFixedSkills)
	teamProtected.Get("/skills/me", skillHandler.GetMySkills)
	teamProtected.Post("/skills/me", skillHandler.SetMySkills)
	teamProtected.Get("/skills/team", skillHandler.GetTeamSkills)
	teamProtected.Put("/skills/member/:id", middleware.RequireRole("captain"), skillHandler.SetMemberSkills)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 GroupTracker API starting on %s", addr)
	log.Fatal(app.Listen(addr))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
