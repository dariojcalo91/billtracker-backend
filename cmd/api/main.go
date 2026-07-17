package main

import (
	"context"
	"log"
	"net/http"
	"os"

	httpadapter "github.com/dariojcalo91/billtracker/internal/adapter/http"
	"github.com/dariojcalo91/billtracker/internal/adapter/postgres"
	"github.com/dariojcalo91/billtracker/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://billtracker:billtracker@localhost:5432/billtracker?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	userRepo := postgres.NewUserRepository(pool)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}
	billRepo := postgres.NewBillRepository(pool)

	registerSvc := usecase.NewRegisterService(userRepo)
	loginSvc := usecase.NewLoginService(userRepo, jwtSecret)
	billSvc := usecase.NewBillService(billRepo)

	authHandler := httpadapter.NewAuthHandler(registerSvc, loginSvc)

	billHandler := httpadapter.NewBillHandler(billSvc)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	protected := router.Group("/")
	protected.Use(httpadapter.AuthMiddleware(jwtSecret))
	{
		protected.POST("/bills", billHandler.Create)
		protected.GET("/bills", billHandler.List)
		protected.GET("/bills/:id", billHandler.Get)
		protected.PUT("/bills/:id", billHandler.Update)
		protected.DELETE("/bills/:id", billHandler.Delete)
	}

	log.Println("server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
