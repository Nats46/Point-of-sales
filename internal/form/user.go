package form

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
    Group    string `json:"group" binding:"required"`
    Spv      *int64 `json:"spv"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
