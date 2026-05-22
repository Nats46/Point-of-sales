package handler

import (
	"fmt"
	"net/http"
	"point-of-sales/internal/form"
	"point-of-sales/internal/model"
	"point-of-sales/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) CreateSales(c *gin.Context) {
	var req form.CreateSalesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	header := &model.SalesHeader{
		CustomerId:    req.CustomerId,
		CashierId:     req.CashierId,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}

	var details []model.SalesDetail
	for _, d := range req.SalesDetails {
		details = append(details, model.SalesDetail{
			ItemId:   d.ItemId,
			ItemName: d.ItemName,
			ItemCode: d.ItemCode,
			Quantity: d.Quantity,
			Price:    d.Price,
			Discount: d.Discount,
			Tax:      d.Tax,
		})
	}

	sales, err := h.svc.CreateSales(c.Request.Context(), header, details, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sales created successfully",
		"data":    sales,
	})
}

func (h *TransactionHandler) GetSales(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sales ID is required"})
		return
	}

	salesID, _ := strconv.ParseInt(id, 10, 64)
	sales, err := h.svc.GetSalesByID(c.Request.Context(), salesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sales retrieved successfully",
		"data":    sales,
	})
}

func (h *TransactionHandler) UpdateSales(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sales ID is required"})
		return
	}

	salesID, _ := strconv.ParseInt(id, 10, 64)

	var req form.UpdateSalesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	header := &model.SalesHeader{
		TransactionId: salesID,
		CustomerId:    req.CustomerId,
		CashierId:     req.CashierId,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}

	var details []model.SalesDetail
	for _, d := range req.SalesDetails {
		details = append(details, model.SalesDetail{
			SalesDetailId: d.SalesDetailId,
			ItemId:        d.ItemId,
			ItemName:      d.ItemName,
			ItemCode:      d.ItemCode,
			Quantity:      d.Quantity,
			Price:         d.Price,
			Discount:      d.Discount,
			Tax:           d.Tax,
		})
	}

	sales, err := h.svc.UpdateSales(c.Request.Context(), header, details, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sales updated successfully",
		"data":    sales,
	})
}

func (h *TransactionHandler) DeleteSales(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sales ID is required"})
		return
	}

	salesID, _ := strconv.ParseInt(id, 10, 64)

	err := h.svc.DeleteSales(c.Request.Context(), salesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sales deleted successfully",
	})
}

func (h *TransactionHandler) ListSales(c *gin.Context) {
	var page, pageSize int64 = 1, 10
	if c.Query("page") != "" {
		fmt.Sscanf(c.Query("page"), "%d", &page)
	}
	if c.Query("pageSize") != "" {
		fmt.Sscanf(c.Query("pageSize"), "%d", &pageSize)
	}
	filter := c.Query("filter")

	sales, err := h.svc.ListSales(c.Request.Context(), page, pageSize, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sales retrieved successfully",
		"data":    sales,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
		},
	})
}