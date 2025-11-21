package testutil

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/gin-gonic/gin"
)

const (
	TenantOne  = "tenant-1"
	UserOne    = "user-1"
	ChannelOne = "channel-1"
)

// SetupTestRouter creates a test Gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// SetupMockDB initializes a mock database for testing
func SetupMockDB(t *testing.T) sqlmock.Sqlmock {
	_, mock, err := db.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to setup mock DB: %v", err)
	}
	return mock
}

// MockUserRows returns mock user data rows
func MockUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "email", "name", "password", "role", "tenant_id"}).
		AddRow(UserOne, "user1@example.com", "User One", "hashedpass", "MEMBER", TenantOne)
}

// MockChannelRows returns mock channel data rows
func MockChannelRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "stream_id", "name", "description", "tenant_id", "created_by"}).
		AddRow(ChannelOne, "stream-123", "General", "General channel", TenantOne, UserOne)
}

// MockChannelMemberRows returns mock channel member data rows
func MockChannelMemberRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "channel_id", "user_id", "tenant_id"}).
		AddRow("member-1", ChannelOne, UserOne, TenantOne)
}

// MockTenantRows returns mock tenant data rows
func MockTenantRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name"}).
		AddRow(TenantOne, "Test Org")
}
