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

	// Load environment variables
	cfg := config.LoadEnv()

	// Database connection
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Auto migration
	if err := db.AutoMigrate(&models.Task{}, &models.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Repositories
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Services (use constructors)
	taskService := service.NewTaskService(taskRepo)
	authService := service.NewAuthService(userRepo)

	// Handlers
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(authService)

	// Background worker context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.StartWorker(ctx, taskRepo, cfg.AutoCompleteMinutes)

	// Gin setup
	r := gin.Default()
	r.Use(middleware.Logger())

	// Swagger route
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

	// HTTP Server for graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server
	go func() {
		log.Println("Server started on port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Stop worker
	cancel()

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
}
