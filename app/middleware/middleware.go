package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingApiGateway/app/models"
	"theAmazingApiGateway/app/security"
)

func ValidateTokenAndPermission(permissionLiteral string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		hasPermissionToAccess, err := hasPermission(token.Id, permissionLiteral)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if hasPermissionToAccess == false {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have enough permissions to access"})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("last_name", token.LastName)
		c.Set("email", token.Email)
	}
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("last_name", token.LastName)
		c.Set("email", token.Email)
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		userData, found, err := models.GetUserById(token.Id)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if userData.RoleID != models.ADMIN {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have enough permissions to access"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(userID uint, permissionLiteral string) (bool, error) {

	//Search if user has user profile permission
	permissionList, err := models.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	var hasPermission = 0
	for _, permissionDescription := range permissionList {
		if permissionDescription == permissionLiteral {
			hasPermission = 1
		}
	}

	if hasPermission == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

func AppendService(serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("serviceURL", serviceURL)
	}
}
