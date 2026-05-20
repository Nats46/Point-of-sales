package handler

import (
    "net/http"
    "point-of-sales/internal/form"
    "point-of-sales/internal/service"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
    var req form.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.svc.Register(c.Request.Context(), req.Username, req.Password, req.Group, req.Spv)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "data":    user,
    })
}

func (h *UserHandler) Login(c *gin.Context) {
    var req form.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}
