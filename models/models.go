package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleModerator Role = "MODERATOR"
	RoleMember    Role = "MEMBER"
	RoleGuest     Role = "GUEST"
)

type Tenant struct {
	ID   string `gorm:"type:uuid;primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;not null" json:"name"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TenantID is set by the backend logic, not by registration request
type User struct {
	ID       string `gorm:"type:uuid;primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Name     string `json:"name"`
	Password string `gorm:"not null" json:"password"`
	Role     Role   `gorm:"default:MEMBER" json:"role"`
	TenantID string `json:"tenant_id"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Channel struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"id"`
	StreamId    string `gorm:"uniqueIndex;not null" json:"stream_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	CreatedBy   string `json:"created_by"`
}

func (c *Channel) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type ChannelMember struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelID string `gorm:"not null;index:idx_channel_user" json:"channel_id"`
	UserID    string `gorm:"not null;index:idx_channel_user" json:"user_id"`
	TenantID  string `gorm:"not null;index" json:"tenant_id"`
	JoinedAt  int64  `gorm:"autoCreateTime" json:"joined_at"`
}

func (cm *ChannelMember) BeforeCreate(tx *gorm.DB) (err error) {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	return nil
}
