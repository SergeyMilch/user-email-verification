package main

import (
	"log"
	"os"

	"github.com/SergeyMilch/user-email-verification/internal/db"
	"github.com/SergeyMilch/user-email-verification/internal/handler"
	"github.com/SergeyMilch/user-email-verification/internal/repository"
	"github.com/SergeyMilch/user-email-verification/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// По умолчанию используем строку подключения для локального запуска
	// Это сделано для удобства локального запуска и тестирования
    // В реальных условиях нельзя хранить чувствительные данные в коде
    // Для продакшен окружения или CI/CD пайплайна нужно задать DATABASE_URL 
    // как переменную окружения, и тогда этот DSN будет автоматически использован
    dsn := "postgres://test-user:password@localhost:5432/test_db?sslmode=disable"
    if val := os.Getenv("DATABASE_URL"); val != "" {
        dsn = val
    }

    pg, err := db.NewPostgresDB(dsn)
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    userRepo := repository.NewUserRepository(pg)
    tokenRepo := repository.NewTokenRepository(pg)
    userService := service.NewUserService(userRepo, tokenRepo)
    userHandler := handler.NewUserHandler(userService)

    r := gin.Default()

    r.POST("/users/register", userHandler.Register)
    r.GET("/users/verify", userHandler.VerifyEmail)

    log.Println("Server started at :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
