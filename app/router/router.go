package router

import (
	"github.com/aviddiviner/gin-limit"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingApiGateway/app"
	"theAmazingApiGateway/app/config"
	"theAmazingApiGateway/app/middleware"
)

var router *gin.Engine

func CreateRouter() {
	router = gin.New()

	router.Use(gin.Logger())
	router.Use(nice.Recovery(recoveryHandler))
	router.Use(limit.MaxAllowed(10))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET,PUT,POST,DELETE"},
		AllowHeaders:    []string{"accept,x-access-token,content-type,authorization"},
	}))

	userAdministration := router.Group("", middleware.AppendService(config.GetConfig().USER_ADMINISTRATION_SERVICE))
	{

		//Oauth2 with Google. Needs frontend
		userAdministration.GET("/", app.CallService)
		userAdministration.GET("/login", app.CallService)
		userAdministration.GET("/googleCallback", app.CallService)

		//Login routes
		userAdministration.POST("/login", app.CallService)
		userAdministration.POST("/signup", app.CallService)
		userAdministration.POST("/recoverPassword", app.CallService)
		userAdministration.PUT("/password", app.CallService)

		//This should be on the front end. When the frontend is made, this should be a POST and the query params would be converted to body params
		userAdministration.GET("/confirmEmail", app.CallService)

		//User management
		userAdministration.GET("/users", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("User Management"), app.CallService)
		userAdministration.PUT("/users/:id", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("User Management"), app.CallService)
		userAdministration.PUT("/users/:id/enable", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("User Management"), app.CallService)

		//Role management
		userAdministration.GET("/roles", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("Role Management"), app.CallService)
		userAdministration.PUT("/roles/:id/permissions", middleware.IsAdmin(), middleware.ValidateTokenAndPermission("Role Management"), app.CallService)

		//User profile endpoints
		userAdministration.PUT("/user/profile", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.GET("/user/profile", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.POST("/user/profile/picture", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.DELETE("/user/profile/picture", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.PUT("/user/password", middleware.ValidateTokenAndPermission("Profile"), app.CallService)

		//Address endpoints
		userAdministration.GET("/user/address", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.POST("/user/address", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.PUT("/user/address/:id", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.PUT("/user/address/:id/main", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.DELETE("/user/address/:id", middleware.ValidateTokenAndPermission("Profile"), app.CallService)

		//Phone endpoints
		userAdministration.PUT("/user/phone", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.POST("/user/phone/verify", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.GET("/user/resendVerificationSMS", middleware.ValidateTokenAndPermission("Profile"), app.CallService)

		//Email change and verification
		userAdministration.PUT("/user/profile/email", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.PUT("/user/verifyEmail", middleware.ValidateTokenAndPermission("Profile"), app.CallService)
		userAdministration.GET("/user/resendConfirmationEmail", middleware.ValidateTokenAndPermission("Profile"), app.CallService)

	}

	postManagement := router.Group("", middleware.AppendService(config.GetConfig().POST_MANAGEMENT_SERVICE))
	{
		//Public
		postManagement.GET("/post", app.CallService)
		postManagement.GET("/post/:id", app.CallService)
		postManagement.GET("/posts/:postID/comment/:id", app.CallService)
		postManagement.GET("/lastPosts", app.CallService)
		postManagement.GET("/lastComments", app.CallService)

		//Post creation
		postManagement.POST("/post", middleware.ValidateToken(), app.CallService)
		postManagement.PUT("/post/:id", middleware.ValidateToken(), app.CallService)
		postManagement.DELETE("/post/:id", middleware.ValidateToken(), app.CallService)
		postManagement.PATCH("/post/:id", middleware.ValidateToken(), app.CallService)

		//Comment creation
		postManagement.POST("/posts/:postID/comment", middleware.ValidateToken(), app.CallService)
		postManagement.PUT("/posts/:postID/comment/:id", middleware.ValidateToken(), app.CallService)
		postManagement.DELETE("/posts/:postID/comment/:id", middleware.ValidateToken(), app.CallService)
		postManagement.PATCH("/posts/:postID/comment/:id", middleware.ValidateToken(), app.CallService)

	}

	notificator := router.Group("", middleware.AppendService(config.GetConfig().NOTIFICATIONS_SERVICE))
	{
		//Notifications
		notificator.GET("/notification", middleware.ValidateToken(), app.CallService)

	}

}

func RunRouter() {
	router.Run(":" + config.GetConfig().PORT)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	detail := ""
	if config.GetConfig().ENV == "develop" {
		detail = err.(error).Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": detail})
}
