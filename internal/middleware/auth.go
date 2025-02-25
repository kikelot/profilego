package middleware

import (
	"net/http"
	"profilego/internal/client"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó un token"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Llamar al microservicio de autenticación para obtener el userId
		user, err := authClient.GetCurrentUser(c.Request.Context(), token)
		if err != nil || user.ID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o usuario no autenticado"})
			c.Abort()
			return
		}

		// Guardar el userId en el contexto de Gin para que los handlers lo usen
		c.Set("userId", user.ID)

		c.Next()
	}
}
