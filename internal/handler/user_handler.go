package handler

import (
	"net/http"

	"github.com/SergeyMilch/user-email-verification/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
    return &UserHandler{service: s}
}

func (h *UserHandler) Register(c *gin.Context) {
    var req struct {
        Nick  string `json:"nick"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.service.RegisterUser(req.Nick, req.Name, req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    verifyLink := "http://localhost:8080/users/verify?token=" + token
    c.JSON(http.StatusOK, gin.H{
        "id":               user.ID,
        "nick":             user.Nick,
        "name":             user.Name,
        "email":            user.Email,
        "email_verified":   user.EmailVerified,
        "verification_link": verifyLink,
        "status":           "pending",
    })
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
        return
    }

    user, err := h.service.VerifyEmail(token)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "verified",
        "email":  user.Email,
    })
}