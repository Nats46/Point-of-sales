package router

import (
    "net/http"

    "point-of-sales/internal/handler"
    "point-of-sales/internal/repository"
    "point-of-sales/internal/service"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB, jwtSecret string) *gin.Engine {
    r := gin.Default()
    r.Use(gin.Logger())

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-API-KEY"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Dependency Injection
    userRepo := repository.NewUserRepository(db)
    userSvc := service.NewUserService(userRepo, jwtSecret)
    userHandler := handler.NewUserHandler(userSvc)

    api := r.Group("/api/v1")
    {
        api.GET("/", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "🎉 Hello from Go! App is running successfully!",
            })
        })

        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })

        // Authentication group
        auth := api.Group("/auth")
        {
            auth.POST("/register", userHandler.Register)
            auth.POST("/login", userHandler.Login)
        }
    }

    return r
}
