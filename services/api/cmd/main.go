package main

//точка входа API: запускает Gin сервер
import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/handler"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file, reading from environment")
	}

	ctx := context.Background()

	pgPool, err := db.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pgPool.Close()
	log.Println("✓ PostgreSQL connected")

	redisClient, err := db.NewRedisClient(ctx, os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("✓ Redis connected")

	// Создаём handler и регистрируем маршруты
	h := handler.New(pgPool, redisClient)

	r := gin.Default()
	h.RegisterRoutes(r)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
