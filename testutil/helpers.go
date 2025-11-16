package testutil

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/gin-gonic/gin"
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
		AddRow("user-1", "user1@example.com", "User One", "hashedpass", "MEMBER", "tenant-1")
}

// MockChannelRows returns mock channel data rows
func MockChannelRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "stream_id", "name", "description", "tenant_id", "created_by"}).
		AddRow("channel-1", "stream-123", "General", "General channel", "tenant-1", "user-1")
}

// MockChannelMemberRows returns mock channel member data rows
func MockChannelMemberRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "channel_id", "user_id", "tenant_id"}).
		AddRow("member-1", "channel-1", "user-1", "tenant-1")
}

// MockTenantRows returns mock tenant data rows
func MockTenantRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name"}).
		AddRow("tenant-1", "Test Org")
}
