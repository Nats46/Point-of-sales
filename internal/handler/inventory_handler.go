package handler

import (
	"fmt"
	"net/http"
	"point-of-sales/internal/form"
	"point-of-sales/internal/service"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	svc service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req form.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.svc.CreateInventory(c.Request.Context(), req.ItemCode, req.ItemName, req.Price, req.Unit, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Inventory created successfully",
		"data":    inventory,
	})
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inventory ID is required"})
		return
	}

	var inventoryID int64
	fmt.Sscanf(id, "%d", &inventoryID)

	inventory, err := h.svc.GetInventory(c.Request.Context(), inventoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory retrieved successfully",
		"data":    inventory,
	})
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inventory ID is required"})
		return
	}

	var inventoryID int64
	fmt.Sscanf(id, "%d", &inventoryID)

	var req form.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.svc.UpdateInventory(c.Request.Context(), inventoryID, req.ItemCode, req.ItemName, req.Price, req.Unit, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory updated successfully",
		"data":    inventory,
	})
}

func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inventory ID is required"})
		return
	}

	var inventoryID int64
	fmt.Sscanf(id, "%d", &inventoryID)

	err := h.svc.DeleteInventory(c.Request.Context(), inventoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory deleted successfully",
	})
}

func (h *InventoryHandler) ListInventories(c *gin.Context) {
	var page, pageSize int64 = 1, 10
	if c.Query("page") != "" {
		fmt.Sscanf(c.Query("page"), "%d", &page)
	}
	if c.Query("pageSize") != "" {
		fmt.Sscanf(c.Query("pageSize"), "%d", &pageSize)
	}
	filter := c.Query("filter")

	inventories, err := h.svc.ListInventories(c.Request.Context(), page, pageSize, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventories retrieved successfully",
		"data":    inventories,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
		},
	})
}