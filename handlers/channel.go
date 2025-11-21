package handlers

import (
	"net/http"

	"log"

	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/Nyagar-Abraham/chat-app/models"
	"github.com/Nyagar-Abraham/chat-app/services"
	"github.com/gin-gonic/gin"
)

// CreateChannel creates a new channel (Admin/Moderator only)
// @Summary Create a channel
// @Description Creates a new chat channel for the tenant and Stream
// @Tags channels
// @Accept json
// @Produce json
// @Param channel body models.Channel true "Channel info"
// @Success 201 {object} models.Channel
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels [post]
func CreateChannel(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != string(models.RoleAdmin) && userRole != string(models.RoleModerator) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions."})
		return
	}
	tenantID, _ := c.Get("tenant_id")
	userId, _ := c.Get("user_id")

	//	create channel
	streamChannelID, err := services.CreateStreamChannel(models.Channel{
		Name:        req.Name,
		Description: req.Description,
		TenantID:    tenantID.(string),
		CreatedBy:   userId.(string),
	}, userId.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Stream Channel"})
		return
	}

	channel := models.Channel{
		StreamId:    streamChannelID,
		Name:        req.Name,
		Description: req.Description,
		TenantID:    tenantID.(string),
		CreatedBy:   userId.(string),
	}

	if err := db.DB.Create(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Stream Channel"})
		return
	}

	member := models.ChannelMember{
		ChannelID: channel.ID,
		UserID:    userId.(string),
		TenantID:  tenantID.(string),
	}
	db.DB.Create(&member)

	c.JSON(http.StatusCreated, channel)
}

// ListChannels lists all channels for a tenant
// @Summary List channels
// @Description Lists all chat channels for the tenant
// @Tags channels
// @Produce json
// @Success 200 {array} models.Channel
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels [get]
func ListChannels(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID not found"})
		return
	}

	var channels []models.Channel
	if err := db.DB.Where("tenant_id = ?", tenantID).Find(&channels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch channels"})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// AddUserToChannel adds a user to a channel (Admin/Moderator only)
// @Summary Add user to channel
// @Description Adds a user to a channel within the same tenant
// @Tags channels
// @Accept json
// @Produce json
// @Param id path string true "Channel ID"
// @Param request body object true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels/{id}/members [post]
func AddUserToChannel(c *gin.Context) {
	log.Println("AddUserToChannel RUN")
	channelID := c.Param("id")
	tenantID, _ := c.Get("tenant_id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("AddUserToChannel BINDERR=", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := services.AddUserToChannel(channelID, req.UserID, tenantID.(string)); err != nil {
		log.Println("AddUserToChannel BINDERR=", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to channel"})
}

// RemoveUserFromChannel removes a user from a channel (Admin/Moderator only)
// @Summary Remove user from channel
// @Description Removes a user from a channel
// @Tags channels
// @Param id path string true "Channel ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels/{id}/members/{user_id} [delete]
func RemoveUserFromChannel(c *gin.Context) {
	channelID := c.Param("id")
	userID := c.Param("user_id")
	tenantID, _ := c.Get("tenant_id")

	if err := services.RemoveUserFromChannel(channelID, userID, tenantID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from channel"})
}

// JoinChannel allows a user to join a channel
// @Summary Join channel
// @Description Allows authenticated user to join a channel in their tenant
// @Tags channels
// @Param id path string true "Channel ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels/{id}/join [post]
func JoinChannel(c *gin.Context) {
	channelID := c.Param("id")
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")

	if err := services.AddUserToChannel(channelID, userID.(string), tenantID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined channel successfully"})
}

// LeaveChannel allows a user to leave a channel
// @Summary Leave channel
// @Description Allows authenticated user to leave a channel
// @Tags channels
// @Param id path string true "Channel ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels/{id}/leave [post]
func LeaveChannel(c *gin.Context) {
	channelID := c.Param("id")
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")

	if err := services.RemoveUserFromChannel(channelID, userID.(string), tenantID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left channel successfully"})
}

// GetChannelMembers lists all members of a channel
// @Summary List channel members
// @Description Lists all users who are members of a channel
// @Tags channels
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /channels/{id}/members [get]
func GetChannelMembers(c *gin.Context) {
	channelID := c.Param("id")
	tenantID, _ := c.Get("tenant_id")

	users, err := services.GetChannelMembers(channelID, tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch members"})
		return
	}
	c.JSON(http.StatusOK, users)
}
