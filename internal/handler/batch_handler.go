package handler

import (
	"fmt"
	"net/http"
	"point-of-sales/internal/form"
	"point-of-sales/internal/service"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	svc service.BatchService
}

func NewBatchHandler(svc service.BatchService) *BatchHandler {
	return &BatchHandler{svc: svc}
}

func (h *BatchHandler) CreateBatch(c *gin.Context) {
	var req form.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch, err := h.svc.CreateBatch(c.Request.Context(), req.ItemCode, req.ItemName, req.BatchNumber, req.Stock, req.Status, req.ExpiredDate, req.StockDate, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Batch created successfully",
		"data":    batch,
	})
}

func (h *BatchHandler) GetBatch(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch ID is required"})
		return
	}

	var batchID int64
	fmt.Sscanf(id, "%d", &batchID)

	batch, err := h.svc.GetBatch(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch retrieved successfully",
		"data":    batch,
	})
}

func (h *BatchHandler) UpdateBatch(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch ID is required"})
		return
	}

	var batchID int64
	fmt.Sscanf(id, "%d", &batchID)

	var req form.UpdateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch, err := h.svc.UpdateBatch(c.Request.Context(), batchID, req.ItemCode, req.ItemName, req.BatchNumber, req.Stock, req.Status, req.ExpiredDate, req.StockDate, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch updated successfully",
		"data":    batch,
	})
}

func (h *BatchHandler) DeleteBatch(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch ID is required"})
		return
	}

	var batchID int64
	fmt.Sscanf(id, "%d", &batchID)

	err := h.svc.DeleteBatch(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch deleted successfully",
	})
}

func (h *BatchHandler) ListBatches(c *gin.Context) {
	var page, pageSize int64 = 1, 10
	if c.Query("page") != "" {
		fmt.Sscanf(c.Query("page"), "%d", &page)
	}
	if c.Query("pageSize") != "" {
		fmt.Sscanf(c.Query("pageSize"), "%d", &pageSize)
	}
	filter := c.Query("filter")

	batches, err := h.svc.ListBatches(c.Request.Context(), page, pageSize, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batches retrieved successfully",
		"data":    batches,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
		},
	})
}