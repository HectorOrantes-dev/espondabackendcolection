package dependencies_login

import (
	"coleccionbackend/src/core/middleware"
	"coleccionbackend/src/feature/login/aplication"
	"coleccionbackend/src/feature/login/infraestructure/adapters"
	"coleccionbackend/src/feature/login/infraestructure/controllers"
	"coleccionbackend/src/feature/login/infraestructure/routers"
	"github.com/gin-gonic/gin"
)

type LoginDependencies struct {
	controller *controllers.LoginController
}

func NewLoginDependencies() *LoginDependencies {
	repo := adapters.NewSupabaseLoginRepository()

	loginUC := aplication.NewLoginUseCase(repo)
	logoutUC := aplication.NewLogoutUseCase(repo)
	refreshUC := aplication.NewRefreshUseCase(repo)

	// Guard compartido para toda la vida del proceso (en memoria)
	guard := middleware.NewBruteForceGuard()

	ctrl := controllers.NewLoginController(loginUC, logoutUC, refreshUC, guard)
	return &LoginDependencies{controller: ctrl}
}

func (d *LoginDependencies) RegisterRoutes(rg *gin.RouterGroup) {
	routers.RegisterLoginRoutes(rg, d.controller)
}
