package main

// @title Task Manager API
// @version 1.0
// @description Task Management Service in Go
// @host localhost:8080
// @BasePath /

import (
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

func main() {
	config.LoadEnv()

	db, _ := gorm.Open(postgres.Open(viper.GetString("DB_URL")), &gorm.Config{})
	db.AutoMigrate(&models.Task{}, &models.User{})

	taskRepo := &repository.TaskRepository{DB: db}
	taskService := &service.TaskService{Repo: taskRepo}

	worker.StartWorker(taskRepo, viper.GetInt("AUTO_COMPLETE_MINUTES"))

	r := gin.Default()

	// routes
	userRepo := &repository.UserRepository{DB: db}
	authService := &service.AuthService{Repo: userRepo}

	authHandler := &handler.AuthHandler{Service: authService}
	taskHandler := &handler.TaskHandler{Service: taskService}

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	auth := r.Group("/tasks")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.POST("", taskHandler.CreateTask)
		auth.GET("", taskHandler.ListTasks)
		auth.GET("/:id", taskHandler.GetTask)
		auth.DELETE("/:id", taskHandler.DeleteTask)
	}

	r.Use(middleware.Logger())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: r,
	// }

	// go func() {
	// 	srv.ListenAndServe()
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)

	// <-quit
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// srv.Shutdown(ctx)

}
