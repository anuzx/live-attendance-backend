package response

import "github.com/gin-gonic/gin"

func ApiResponse[T any](c *gin.Context, status int, data T) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func ApiError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}
