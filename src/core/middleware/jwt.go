package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwksCache mantiene el resolutor de claves públicas (JWKS) de Supabase.
// Se inicializa una sola vez y refresca las claves en segundo plano.
var (
	jwksOnce sync.Once
	jwksKf   keyfunc.Keyfunc
	jwksErr  error
)

func getJWKS() (keyfunc.Keyfunc, error) {
	jwksOnce.Do(func() {
		base := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
		url := base + "/auth/v1/.well-known/jwks.json"
		jwksKf, jwksErr = keyfunc.NewDefaultCtx(context.Background(), []string{url})
	})
	return jwksKf, jwksErr
}

// keyResolver elige la clave de verificación según el algoritmo del token:
//   - HS256 (legacy JWT secret) → usa JWT_SECRET del .env
//   - ES256/RS256 (nuevas JWT Signing Keys) → usa la clave pública del JWKS
func keyResolver(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}
	kf, err := getJWKS()
	if err != nil {
		return nil, err
	}
	return kf.Keyfunc(t)
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}

		token, err := jwt.Parse(parts[1], keyResolver)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "claims inválidos"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("email", claims["email"])
		c.Next()
	}
}
