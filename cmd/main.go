package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/flagman-backend/internal/config"
	"github.com/immxrtalbeast/flagman-backend/internal/config/controller"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/middleware"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/document"
	"github.com/immxrtalbeast/flagman-backend/internal/usecase/user"
	"github.com/immxrtalbeast/flagman-backend/storage/supabase"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting application", slog.Any("config", cfg))
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	host := os.Getenv("host")
	userDB := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	port := os.Getenv("port")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require pgbouncer=true connect_timeout=10 pool_mode=transaction statement_cache_mode=describe",
		host, userDB, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            false, // Отключаем подготовленные выражения
		SkipDefaultTransaction: true,
		ConnPool:               nil,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{}, &domain.Organization{}, &domain.Department{}, &domain.Document{}, &domain.DocumentRecipient{})
	if err := db.Exec("DEALLOCATE ALL").Error; err != nil {
		panic(err)
	}

	usrRepo := supabase.NewUserRepository(db)
	userINT := user.NewUserInteractor(usrRepo, cfg.TokenTTL, cfg.AppSecret)
	userController := controller.NewUserController(userINT, cfg.TokenTTL, cfg.AppSecret)

	documentRepo := supabase.NewDocumentRepository(db)
	documentINT := document.NewDocumentInteractor(documentRepo, usrRepo)
	documentController := controller.NewDocumentController(documentINT)

	authMiddleware := middleware.AuthMiddleware(cfg.AppSecret)

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
		organization := api.Group("/organization")
		organization.Use(authMiddleware)
		{
			organization.GET("/", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})
		}
		document := api.Group("/document")
		document.Use(authMiddleware)
		{
			document.POST("/create", documentController.CreateDocument)
		}
		user := api.Group("/user")
		user.Use(authMiddleware)
		{
			user.GET("/:id", userController.User)
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
