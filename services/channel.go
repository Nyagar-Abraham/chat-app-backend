package services

import (
	"context"
	"errors"

	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/Nyagar-Abraham/chat-app/models"
)

func AddUserToChannel(channelID, userID, tenantID string) error {
	var channel models.Channel
	if err := db.DB.Where("id = ?::uuid AND tenant_id = ?", channelID, tenantID).First(&channel).Error; err != nil {
		return errors.New("channel not found or access denied")
	}

	var user models.User
	if err := db.DB.Where("id = ?::uuid AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		return errors.New("user not found or access denied")
	}

	var existing models.ChannelMember
	if err := db.DB.Where("channel_id = ? AND user_id = ?", channelID, userID).First(&existing).Error; err == nil {
		return errors.New("user already in channel")
	}

	member := models.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
		TenantID:  tenantID,
	}

	if err := db.DB.Create(&member).Error; err != nil {
		return err
	}

	client := GetStreamClient()
	ch := client.Channel("messaging", channel.StreamId)
	_, err := ch.AddMembers(context.Background(), []string{userID})
	if err != nil {
		db.DB.Delete(&member)
		return errors.New("failed to add user to stream channel: " + err.Error())
	}

	return nil
}

func RemoveUserFromChannel(channelID, userID, tenantID string) error {
	var channel models.Channel
	if err := db.DB.Where("id = ?::uuid AND tenant_id = ?", channelID, tenantID).First(&channel).Error; err != nil {
		return errors.New("channel not found or access denied")
	}

	if err := db.DB.Where("channel_id = ? AND user_id = ? AND tenant_id = ?", channelID, userID, tenantID).Delete(&models.ChannelMember{}).Error; err != nil {
		return err
	}

	client := GetStreamClient()
	ch := client.Channel("messaging", channel.StreamId)
	_, err := ch.RemoveMembers(context.Background(), []string{userID}, nil)
	if err != nil {
		return errors.New("failed to remove user from stream channel: " + err.Error())
	}

	return nil
}

func IsUserChannelMember(channelID, userID, tenantID string) bool {
	var count int64
	db.DB.Model(&models.ChannelMember{}).Where("channel_id = ? AND user_id = ? AND tenant_id = ?", channelID, userID, tenantID).Count(&count)
	return count > 0
}

func GetChannelMembers(channelID, tenantID string) ([]models.User, error) {
	var users []models.User
	err := db.DB.Table("users").
		Joins("JOIN channel_members ON users.id::text = channel_members.user_id").
		Where("channel_members.channel_id = ? AND channel_members.tenant_id = ?", channelID, tenantID).
		Find(&users).Error
	return users, err
}
