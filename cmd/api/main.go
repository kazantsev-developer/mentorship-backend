package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/config"
	"github.com/kazantsev/mentorship-backend/internal/handlers"
	"github.com/kazantsev/mentorship-backend/internal/middleware"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"github.com/kazantsev/mentorship-backend/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := db.InitDB(cfg); err != nil {
		log.Fatal("Failed to init database:", err)
	}

	userRepo := repositories.NewUserRepository(db.GetDB())
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	blockRepo := repositories.NewBlockRepository(db.GetDB())
	roadmapService := services.NewRoadmapService(blockRepo, db.GetDB())
	roadmapHandler := handlers.NewRoadmapHandler(roadmapService)

	materialRepo := repositories.NewMaterialRepository(db.GetDB())
	progressRepo := repositories.NewProgressRepository(db.GetDB())

	bonusRepo := repositories.NewBonusRepository(db.GetDB())
	bonusService := services.NewBonusService(bonusRepo)
	achievementRepo := repositories.NewAchievementRepository(db.GetDB())
	achievementService := services.NewAchievementService(achievementRepo, bonusService)

	progressService := services.NewProgressService(progressRepo, materialRepo, achievementService)
	progressHandler := handlers.NewProgressHandler(progressService)

	assignmentService := services.NewAssignmentService(db.GetDB(), userRepo)
	assignmentHandler := handlers.NewAssignmentHandler(assignmentService)

	interviewService := services.NewInterviewService(db.GetDB())
	interviewHandler := handlers.NewInterviewHandler(interviewService)

	calendarService := services.NewCalendarService(db.GetDB())
	calendarHandler := handlers.NewCalendarHandler(calendarService)

	oneOnOneService := services.NewOneOnOneService(db.GetDB(), bonusService)
	oneOnOneHandler := handlers.NewOneOnOneHandler(oneOnOneService)

	finalCheckService := services.NewFinalCheckService(db.GetDB())
	finalCheckHandler := handlers.NewFinalCheckHandler(finalCheckService)

	profileHandler := handlers.NewProfileHandler(userRepo)

	blockApproveHandler := handlers.NewBlockApproveHandler(progressRepo)
	adminUserHandler := handlers.NewAdminUserHandler(userRepo, authService)
	adminRoadmapHandler := handlers.NewAdminRoadmapHandler(db.GetDB())

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/user/profile", func(c *gin.Context) {
			userID := c.GetString("userID")
			user, err := userRepo.FindByID(userID)
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to fetch user"})
				return
			}
			user.PasswordHash = ""
			c.JSON(200, user)
		})
		protected.PUT("/user/profile", profileHandler.UpdateProfile)
		protected.GET("/user/:user_id/profile", profileHandler.PublicProfile)

		protected.GET("/roadmap", roadmapHandler.GetRoadmap)
		protected.POST("/materials/view", progressHandler.MarkMaterialViewed)

		protected.GET("/bonus/balance", func(c *gin.Context) {
			userID := c.GetString("userID")
			bal, err := bonusService.GetBalance(userID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"balance": bal})
		})
		protected.GET("/bonus/history", func(c *gin.Context) {
			userID := c.GetString("userID")
			history, err := bonusService.GetHistory(userID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, history)
		})
		protected.POST("/bonus/convert", func(c *gin.Context) {
			var req struct {
				BonusAmount int `json:"bonus_amount"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetString("userID")
			discount, err := bonusService.ConvertBonusToDiscount(userID, req.BonusAmount)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"discount_percent": discount})
		})
		protected.GET("/achievements", func(c *gin.Context) {
			achievements, err := achievementRepo.GetActiveAchievements()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetString("userID")
			result := make([]map[string]interface{}, 0)
			for _, ach := range achievements {
				has, _ := achievementRepo.HasUserAchievement(userID, ach.ID)
				result = append(result, map[string]interface{}{
					"id":           ach.ID,
					"title":        ach.Title,
					"description":  ach.Description,
					"reward_bonus": ach.RewardBonus,
					"image_url":    ach.ImageURL,
					"unlocked":     has,
				})
			}
			c.JSON(200, result)
		})

		protected.POST("/blocks/approve", blockApproveHandler.ApproveBlock)

		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.POST("/assign-buddy", assignmentHandler.AssignBuddy)

			adminGroup.GET("/users", adminUserHandler.ListUsers)
			adminGroup.POST("/users", adminUserHandler.CreateUser)
			adminGroup.PUT("/users/:user_id", adminUserHandler.UpdateUser)
			adminGroup.DELETE("/users/:user_id", adminUserHandler.DeleteUser)
			adminGroup.POST("/users/:user_id/change-password", adminUserHandler.ChangePassword)

			adminGroup.GET("/blocks", adminRoadmapHandler.ListBlocks)
			adminGroup.POST("/blocks", adminRoadmapHandler.CreateBlock)
			adminGroup.PUT("/blocks/:id", adminRoadmapHandler.UpdateBlock)
			adminGroup.DELETE("/blocks/:id", adminRoadmapHandler.DeleteBlock)
		}

		protected.GET("/my-students", assignmentHandler.MyStudents)

		protected.POST("/interviews/real", interviewHandler.CreateReal)
		protected.POST("/interviews/mock", interviewHandler.CreateMock)
		protected.POST("/interviews/mock/complete", interviewHandler.CompleteMock)
		protected.GET("/interviews/my", interviewHandler.MyInterviews)
		protected.GET("/interviews/real", interviewHandler.AllReal)

		protected.POST("/calendar/events", calendarHandler.CreateEvent)
		protected.GET("/calendar/events", calendarHandler.MyEvents)
		protected.GET("/calendar/upcoming", calendarHandler.UpcomingEvents)

		protected.POST("/one-on-one", oneOnOneHandler.CreateRequest)
		protected.GET("/one-on-one", oneOnOneHandler.ListRequests)
		protected.POST("/one-on-one/approve", oneOnOneHandler.Approve)
		protected.POST("/one-on-one/reject", oneOnOneHandler.Reject)

		protected.POST("/final-checks/schedule", finalCheckHandler.Schedule)
		protected.POST("/final-checks/complete", finalCheckHandler.Complete)
		protected.GET("/final-checks/student/:student_id", finalCheckHandler.GetForStudent)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
