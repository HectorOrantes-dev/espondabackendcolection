package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/core/middleware"
	"coleccionbackend/src/feature/login/aplication"
	"coleccionbackend/src/feature/login/domain/entities"
)

type LoginController struct {
	loginUC   *aplication.LoginUseCase
	logoutUC  *aplication.LogoutUseCase
	refreshUC *aplication.RefreshUseCase
	guard     *middleware.BruteForceGuard
}

func NewLoginController(
	l *aplication.LoginUseCase,
	lo *aplication.LogoutUseCase,
	r *aplication.RefreshUseCase,
	guard *middleware.BruteForceGuard,
) *LoginController {
	return &LoginController{loginUC: l, logoutUC: lo, refreshUC: r, guard: guard}
}

// Login godoc
// @Summary Iniciar sesión
// @Tags auth
// @Accept json
// @Produce json
// @Param body body entities.LoginRequest true "Credenciales"
// @Success 200 {object} entities.TokenResponse
// @Failure 401 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /auth/login [post]
func (ctrl *LoginController) Login(c *gin.Context) {
	var req entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si la cuenta está bloqueada antes de intentar
	if err := ctrl.guard.CheckBlocked(req.Email); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := ctrl.loginUC.Execute(c.Request.Context(), req)
	if err != nil {
		// Registrar el fallo y aplicar bloqueo si corresponde
		blocked, duration := ctrl.guard.RecordFailure(req.Email)
		attemptsLeft := ctrl.guard.AttemptsLeft(req.Email)

		if blocked {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":    buildBlockMessage(duration),
				"blocked":  true,
				"duration": duration.String(),
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "credenciales incorrectas",
			"attempts_left": attemptsLeft,
		})
		return
	}

	// Login exitoso — reiniciar el contador
	ctrl.guard.RecordSuccess(req.Email)
	c.JSON(http.StatusOK, token)
}

// Logout godoc
// @Summary Cerrar sesión
// @Tags auth
// @Security BearerAuth
// @Success 204
// @Router /auth/logout [post]
func (ctrl *LoginController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token requerido"})
		return
	}

	if err := ctrl.logoutUC.Execute(c.Request.Context(), parts[1]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Refresh godoc
// @Summary Renovar token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body entities.RefreshRequest true "Refresh token"
// @Success 200 {object} entities.TokenResponse
// @Router /auth/refresh [post]
func (ctrl *LoginController) Refresh(c *gin.Context) {
	var req entities.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.refreshUC.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

func buildBlockMessage(d interface{ String() string }) string {
	switch d.String() {
	case "5m0s":
		return "demasiados intentos fallidos, cuenta bloqueada por 5 minutos"
	case "20m0s":
		return "demasiados intentos fallidos, cuenta bloqueada por 20 minutos"
	default:
		return "demasiados intentos fallidos, cuenta bloqueada temporalmente"
	}
}
