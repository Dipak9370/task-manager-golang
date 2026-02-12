package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "task-manager/confic"
	"task-manager/handler"
	"task-manager/middleware"
	"task-manager/models"
	"task-manager/repository"
	"task-manager/service"
	"task-manager/worker"

	_ "task-manager/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Task Manager API
// @version 1.0
// @description Task Management Service in Go
// @host localhost:8080
// @BasePath /

func main() {
	// Load env
	cfg := config.LoadEnv()

	// DB connection
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.Task{}, &models.User{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Repositories
	taskRepo := &repository.TaskRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}

	// Services
	taskService := &service.TaskService{Repo: taskRepo}
	authService := &service.AuthService{Repo: userRepo}

	// Handlers
	taskHandler := &handler.TaskHandler{Service: taskService}
	authHandler := &handler.AuthHandler{Service: authService}

	// Start background worker with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.StartWorker(ctx, taskRepo, cfg.AutoCompleteMinutes)

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Protected routes
	auth := r.Group("/tasks")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.POST("", taskHandler.CreateTask)
		auth.GET("", taskHandler.ListTasks)
		auth.GET("/:id", taskHandler.GetTask)
		auth.DELETE("/:id", taskHandler.DeleteTask)
	}

	port := viper.GetString("PORT")

	// HTTP server (for graceful shutdown)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Run server in goroutine
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Println("Server started on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	cancel() // stop worker

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
}
