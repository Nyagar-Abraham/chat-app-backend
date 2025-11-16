package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/Nyagar-Abraham/chat-app/models"
	"github.com/Nyagar-Abraham/chat-app/services"
	"github.com/Nyagar-Abraham/chat-app/utils"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=ADMIN MODERATOR MEMBER GUEST"`
	OrgName  string `json:"org_name" binding:"required"`
}

type RegisterResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	TenantId string `json:"tenant_id"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantId string `json:"tenant_id"`
}

// Login authenticates a user and returns a JWT token
// @Summary Login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	//	validate the user (fetch from db, check password)
	var user models.User
	if err := db.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := services.CheckPassword(request.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	//	Issue JWT token with all required claims
	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		Message: user.Name,
	})

}

// Register handles user registration (sign up)
// @Summary Register a new user
// @Description Register a new user (admin, moderator, or user)
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Registration info"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {

	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// validate role
	role := strings.ToUpper(request.Role)
	if role != "ADMIN" && role != "MODERATOR" && role != "GUEST" && role != "MEMBER" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role: must be ADMIN, MODERATOR, MEMBER, or GUEST"})
		return
	}

	// create or find tenant
	var tenant models.Tenant
	if request.OrgName != "" {
		var err error
		tenant, err = findOrCreateTenant(request.OrgName)
		if err != nil {
			log.Printf("Error with tenant: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not handle tenant"})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization name is required"})
		return
	}

	hash, err := services.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hash),
		Role:     models.Role(role),
		TenantID: tenant.ID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	if err := services.CreateStreamUser(user); err != nil {
		log.Printf("Failed to create stream user for %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create stream user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     string(user.Role),
		Token:    token,
		TenantId: user.TenantID,
	})
}

func GetCurrentUser(context *gin.Context) {
	user := models.User{}
	userId, _ := context.Get("user_id")

	if err := db.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}

	context.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     string(user.Role),
		TenantId: user.TenantID,
	})
}

// findOrCreateTenant finds an existing tenant by name or creates a new one
func findOrCreateTenant(orgName string) (models.Tenant, error) {
	var tenant models.Tenant

	// Validate input
	if strings.TrimSpace(orgName) == "" {
		return tenant, fmt.Errorf("organization name cannot be empty")
	}

	// Use GORM's FirstOrCreate - atomic operation that finds or creates
	result := db.DB.Where(models.Tenant{Name: orgName}).FirstOrCreate(&tenant)
	if result.Error != nil {
		log.Printf("Error in FirstOrCreate for tenant '%s': %v", orgName, result.Error)
		return tenant, result.Error
	}

	// Check if this was a new creation
	if result.RowsAffected > 0 {
		log.Printf("Created new tenant: '%s' (ID: %s)", tenant.Name, tenant.ID)
	} else {
		log.Printf("Found existing tenant: '%s' (ID: %s)", tenant.Name, tenant.ID)
	}

	return tenant, nil
}
