package handler

import (
	"fmt"
	"net/http"
	"point-of-sales/internal/form"
	"point-of-sales/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	svc service.DiscountService
}

func NewDiscountHandler(svc service.DiscountService) *DiscountHandler {
	return &DiscountHandler{svc: svc}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var req form.CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := h.svc.CreateDiscount(c.Request.Context(), req.DiscountCode, req.DiscountType, req.Value, req.StartDate, req.EndDate, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Discount created successfully",
		"data":    discount,
	})
}

func (h *DiscountHandler) GetDiscount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discount ID is required"})
		return
	}

	discountID, _ := strconv.ParseInt(id, 10, 64)
	discount, err := h.svc.GetDiscountByID(c.Request.Context(), discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discount retrieved successfully",
		"data":    discount,
	})
}

func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discount ID is required"})
		return
	}

	discountID, _ := strconv.ParseInt(id, 10, 64)

	var req form.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := h.svc.UpdateDiscount(c.Request.Context(), discountID, req.DiscountCode, req.DiscountType, req.Value, req.StartDate, req.EndDate, req.IsActive, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discount updated successfully",
		"data":    discount,
	})
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discount ID is required"})
		return
	}

	discountID, _ := strconv.ParseInt(id, 10, 64)

	err := h.svc.DeleteDiscount(c.Request.Context(), discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discount deleted successfully",
	})
}

func (h *DiscountHandler) ListDiscounts(c *gin.Context) {
	var page, pageSize int64 = 1, 10
	if c.Query("page") != "" {
		fmt.Sscanf(c.Query("page"), "%d", &page)
	}
	if c.Query("pageSize") != "" {
		fmt.Sscanf(c.Query("pageSize"), "%d", &pageSize)
	}
	filter := c.Query("filter")

	discounts, err := h.svc.ListDiscounts(c.Request.Context(), page, pageSize, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discounts retrieved successfully",
		"data":    discounts,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
		},
	})
}