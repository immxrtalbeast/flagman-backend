package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/config"
	"github.com/immxrtalbeast/flagman-backend/internal/controller"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/middleware"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/document"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/enterprise"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/notifications"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/recipient"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/user"
	"github.com/immxrtalbeast/flagman-backend/storage/supabase"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	//TODO: заменить ручку users на юзеров из своей орги
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting application", slog.Any("config", cfg))
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	password := os.Getenv("password")

	dsn := fmt.Sprintf("postgresql://postgres.nrcxwtnzqivegpqkjtue:%s@aws-0-eu-north-1.pooler.supabase.com:5432/postgres", password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{}, &domain.Invitation{}, &domain.Enterprise{}, &domain.Document{}, &domain.DocumentRecipient{})
	if err := db.Exec("DEALLOCATE ALL").Error; err != nil {
		panic(err)
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis" // значение по умолчанию
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379" // значение по умолчанию
	}

	var redisClient *redis.Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // Если используете пароль
		DB:       0,
	})
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}
	fmt.Println("Redis connected:", pong)
	usrRepo := supabase.NewUserRepository(db)
	userINT := user.NewUserInteractor(usrRepo, cfg.TokenTTL, cfg.AppSecret)
	userController := controller.NewUserController(userINT, cfg.TokenTTL, cfg.AppSecret, redisClient)

	notifRepo := supabase.NewNotificationRepository(db)
	enterpriseRepo := supabase.NewEnterpriseRepository(db)
	notifINT := notifications.NewNotificationInteractor(notifRepo, enterpriseRepo)
	notifController := controller.NewNotificationController(notifINT)

	enterpriseINT := enterprise.NewEnterpriseInteractor(enterpriseRepo, usrRepo, notifRepo)
	enterpriseController := controller.NewEnterpriseController(enterpriseINT)

	recipientRepo := supabase.NewDocumentRecipientRepository(db)
	recipientINT := recipient.NewDocumentRecipientInteractor(recipientRepo, usrRepo, os.Getenv("secret_salt"))
	recipientController := controller.NewRecipientController(recipientINT, recipientRepo, redisClient, usrRepo, os.Getenv("EmailService"))

	documentRepo := supabase.NewDocumentRepository(db)
	documentINT := document.NewDocumentInteractor(documentRepo, usrRepo)
	documentController := controller.NewDocumentController(documentINT, recipientRepo)

	authMiddleware := middleware.AuthMiddleware(cfg.AppSecret, redisClient)

	// Настройка маршрутов
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{
		"Authorization",
		"Content-Type",
		"Origin",
		"Accept",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	router.Use(cors.New(config))
	api := router.Group("/api/v1")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)
		api.POST("/logout", userController.Logout).Use(authMiddleware)

		enterprise := api.Group("/enterprise")
		enterprise.Use(authMiddleware)
		{
			enterprise.POST("/create", enterpriseController.CreateEnterprise)
			enterprise.GET("/:id", enterpriseController.Enterprise)
			enterprise.GET("/my", enterpriseController.EnterprisesByUserID)
			enterprise.POST("/invite", enterpriseController.InviteUser)
		}
		notification := api.Group("/notification")
		notification.Use(authMiddleware)
		{
			notification.GET("/my", notifController.MyNotifications)
			notification.POST("/accept", notifController.AcceptInvite)
		}
		document := api.Group("/document")
		document.Use(authMiddleware)
		{
			document.POST("/create", documentController.CreateDocument)
			document.GET("/list", recipientController.ListUserDocuments)
			document.POST("/reject/:id", recipientController.RejectDocument)
			document.POST("/sign/:id", recipientController.SignDocument)
			document.POST("/sign/:id/request", recipientController.RequestToSign)
			document.GET("/:id", recipientController.ByID)

		}
		user := api.Group("/user")
		user.Use(authMiddleware)
		{
			// user.GET("/list", userController.Users)
			user.GET("/:id", userController.User)
			user.GET("/all", userController.UsersAll)
		}
	}
	router.Run(":8080")
}
func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return log
}
