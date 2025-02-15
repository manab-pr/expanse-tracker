package middlewares

import (
	helper "expanse-tracker/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if (clientToken) == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No authorization token provided.")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return

		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("user_id", claims.User_Id)

		c.Next()

	}

}
