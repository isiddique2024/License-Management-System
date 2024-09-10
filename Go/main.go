package main

import (
	"context"
	"fmt"
	"os"

	_ "backend/docs"
	"backend/internal/controllers"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/fasthttp/router"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ctx = context.Background()

func main() {
	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Test Redis connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Could not connect to Redis: %v", err))
	}
	fmt.Println("Connected to Redis")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "data/test.db"
	}

	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		PrepareStmt: true, // Prepared statements enabled
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Application{}, &models.License{})

	// Create a new Gin router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware(), middleware.SecurityHeadersMiddleware(), middleware.CSPMiddleware())

	// Inject Redis client into the Gin context
	r.Use(func(c *gin.Context) {
		c.Set("redisClient", redisClient)
		c.Next()
	})

	controllers.RegisterRoutes(r, db)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve gin router via fasthttp
	fasthttpRouter := router.New()
	fasthttpRouter.NotFound = fasthttpadaptor.NewFastHTTPHandler(r)

	if err := fasthttp.ListenAndServe(":8001", fasthttpRouter.Handler); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
