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

	inventoryRepo := repository.NewInventoryRepository(db)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventorySvc)

	batchRepo := repository.NewBatchRepository(db)
	batchSvc := service.NewBatchService(batchRepo)
	batchHandler := handler.NewBatchHandler(batchSvc)

	discountRepo := repository.NewDiscountRepository(db)
	discountSvc := service.NewDiscountService(discountRepo)
	discountHandler := handler.NewDiscountHandler(discountSvc)

	transactionRepo := repository.NewTransactionRepository(db)
	detailRepo := repository.NewSalesDetailRepository(db)
	transactionSvc := service.NewTransactionService(transactionRepo, detailRepo)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

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

		// Inventory group
		inventory := api.Group("/inventories")
		{
			inventory.POST("", inventoryHandler.CreateInventory)
			inventory.GET("/:id", inventoryHandler.GetInventory)
			inventory.PUT("/:id", inventoryHandler.UpdateInventory)
			inventory.DELETE("/:id", inventoryHandler.DeleteInventory)
			inventory.GET("", inventoryHandler.ListInventories)
		}

		// Batch group
		batch := api.Group("/batches")
		{
			batch.POST("", batchHandler.CreateBatch)
			batch.GET("/:id", batchHandler.GetBatch)
			batch.PUT("/:id", batchHandler.UpdateBatch)
			batch.DELETE("/:id", batchHandler.DeleteBatch)
			batch.GET("", batchHandler.ListBatches)
		}

		// Discount group
		discount := api.Group("/discounts")
		{
			discount.POST("", discountHandler.CreateDiscount)
			discount.GET("/:id", discountHandler.GetDiscount)
			discount.PUT("/:id", discountHandler.UpdateDiscount)
			discount.DELETE("/:id", discountHandler.DeleteDiscount)
			discount.GET("", discountHandler.ListDiscounts)
		}

		// Transaction group
		transaction := api.Group("/transactions")
		{
			transaction.POST("", transactionHandler.CreateSales)
			transaction.GET("/:id", transactionHandler.GetSales)
			transaction.PUT("/:id", transactionHandler.UpdateSales)
			transaction.DELETE("/:id", transactionHandler.DeleteSales)
			transaction.GET("", transactionHandler.ListSales)
		}
	}

	return r
}
