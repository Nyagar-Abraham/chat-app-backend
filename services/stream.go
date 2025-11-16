package services

import (
	"context"
	"errors"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/Nyagar-Abraham/chat-app/models"
	"github.com/google/uuid"

	"log"
	"os"
	"sync"
	"time"
)

var streamClient *stream.Client
var once sync.Once

// return a singleton stream client
func GetStreamClient() *stream.Client {
	once.Do(func() {
		apiKey := os.Getenv("STREAM_API_KEY")
		apiSecret := os.Getenv("STREAM_API_SECRET")
		if apiKey == "" || apiSecret == "" {
			log.Fatal("STREAM_API_KEY and STREAM_API_SECRET must be set")
		}
		var err error
		streamClient, err = stream.NewClient(apiKey, apiSecret)
		if err != nil {
			log.Fatalf("Failed to create stream client: %v", err)
		}
	})
	return streamClient
}

// generate a stream chat token for a user
func CreateStreamToken(userId string) (string, error) {
	if userId == "" {
		return "", errors.New("userId must is required for token generation")
	}
	c := GetStreamClient()
	token, err := c.CreateToken(userId, time.Now().Add(24*time.Hour))
	return token, err
}

func mapRoleToStream(role models.Role) string {
	switch role {
	case models.RoleAdmin:
		return "admin"
	case models.RoleModerator:
		return "channel_moderator"
	case models.RoleMember:
		return "user"
	case models.RoleGuest:
		return "guest"
	default:
		return "user"
	}
}

func CreateStreamUser(user models.User) error {
	client := GetStreamClient()
	_, err := client.UpsertUsers(context.Background(), &stream.User{
		ID:   user.ID,
		Name: user.Name,
		Role: mapRoleToStream(user.Role),
		ExtraData: map[string]interface{}{
			"tenant_id": user.TenantID,
			"email":     user.Email,
			"app_role":  string(user.Role),
		},
	})
	if err != nil {
		log.Printf("CreateStreamUser error for user %s: %v", user.ID, err)
	}
	return err
}

// create a channel
func CreateStreamChannel(channel models.Channel, creatorID string) (string, error) {
	client := GetStreamClient()
	//	Ensure channelId is less than 64 characters for stream
	shortTenantID := channel.TenantID

	if len(shortTenantID) > 8 {
		shortTenantID = shortTenantID[:8]
	}

	channelID := shortTenantID + "-" + uuid.New().String()
	ch, err := client.CreateChannel(
		context.Background(),
		"messaging",
		channelID,
		creatorID,
		&stream.ChannelRequest{
			Members: []string{creatorID},
			ExtraData: map[string]interface{}{
				"tenant_id":   channel.TenantID,
				"name":        channel.Name,
				"description": channel.Description,
			},
		},
	)
	if err != nil {
		log.Printf("Stream CreateChannel error: %v", err)
		return "", err
	}

	return ch.Channel.ID, nil
}
