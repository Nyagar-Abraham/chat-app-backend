package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Nyagar-Abraham/chat-app/db"
	"github.com/Nyagar-Abraham/chat-app/handlers"
	"github.com/Nyagar-Abraham/chat-app/middleware"
	"github.com/Nyagar-Abraham/chat-app/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on environment variables")
	}
	//	connect db
	db.Connect()

	//	Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	router := gin.Default()
	router.Use(cors.New(config))

	//	swagger docs endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the chat API!!",
		})
	})

	//	Auth endpoints
	router.POST("/auth/login", handlers.Login)
	router.POST("/auth/register", handlers.Register)
	router.GET("/me", middleware.JWTAuth(), handlers.GetCurrentUser)

	//steam Chat token endpoint
	router.GET("/stream/token", middleware.JWTAuth(), handlers.StreamToken)

	//Tenant endpoints
	router.POST("/tenants", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin)), handlers.CreateTenant)
	router.GET("/tenants", handlers.ListTenants)
	router.GET("/tenants/:id", handlers.GetTenant)

	// User endpoints (Admin/Moderator for create/update, Admin for delete, all roles for list)
	router.POST("/users", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin), string(models.RoleModerator)), handlers.CreateUser)
	router.GET("/users", middleware.JWTAuth(), handlers.ListUsers)
	router.PUT("/users/:id", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin), string(models.RoleModerator)), handlers.UpdateUser)
	router.DELETE("/users/:id", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin)), handlers.DeleteUser)

	// Channel endpoints (Admin/Moderator for create, all roles for list)
	router.POST("/channels", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin), string(models.RoleModerator)), handlers.CreateChannel)
	router.GET("/channels", middleware.JWTAuth(), handlers.ListChannels)

	// Channel membership endpoints
	router.POST("/channels/:id/members", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin), string(models.RoleModerator)), handlers.AddUserToChannel)
	router.DELETE("/channels/:id/members/:user_id", middleware.JWTAuth(), middleware.RequireRole(string(models.RoleAdmin), string(models.RoleModerator)), handlers.RemoveUserFromChannel)
	router.GET("/channels/:id/members", middleware.JWTAuth(), handlers.GetChannelMembers)
	router.POST("/channels/:id/join", middleware.JWTAuth(), handlers.JoinChannel)
	router.POST("/channels/:id/leave", middleware.JWTAuth(), handlers.LeaveChannel)

	// Messages endpoint (all authenticated users)
	router.POST("/messages", middleware.JWTAuth(), handlers.SendMessage)
	router.GET("/messages/:stream_id", middleware.JWTAuth(), handlers.GetMessages)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	log.Printf("Listening on port: %s", port)

	if err := router.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
