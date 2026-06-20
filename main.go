package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"coleccionbackend/src/core/database"
	dependencies_coleccion "coleccionbackend/src/feature/coleccion/infraestructure/dependencies_coleccion"
	dependencies_etiquetas "coleccionbackend/src/feature/etiquetas/infraestructure/dependencies_etiquetas"
	dependencies_login "coleccionbackend/src/feature/login/infraestructure/dependencies_login"
	_ "coleccionbackend/docs"
)

// @title           Coleccion Backend API
// @version         1.0
// @description     API para gestión de vehículos de colección con almacenamiento de imágenes en Google Drive.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontró archivo .env, usando variables de entorno del sistema")
	}

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close()

	imgService, err := database.NewGoogleDriveService()
	if err != nil {
		log.Fatal("Error inicializando Google Drive:", err)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20

	// Confiar solo en el proxy local (loopback). Evita que se falsifique
	// la IP de origen vía X-Forwarded-For. Ajusta si despliegas detrás de
	// un proxy/CDN concreto (ej: r.SetTrustedProxies([]string{"10.0.0.1"})).
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS: en producción define ALLOWED_ORIGINS (ej: "https://miapp.vercel.app")
	// separado por comas. Si no está definido, permite todos (útil en local).
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: false,
	}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		corsConfig.AllowOrigins = strings.Split(origins, ",")
	} else {
		corsConfig.AllowAllOrigins = true
	}
	r.Use(cors.New(corsConfig))

	// Health check para Render / monitoreo
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	loginDeps := dependencies_login.NewLoginDependencies()
	loginDeps.RegisterRoutes(api)

	coleccionDeps := dependencies_coleccion.NewColeccionDependencies(db, imgService)
	coleccionDeps.RegisterRoutes(api)

	etiquetasDeps := dependencies_etiquetas.NewEtiquetasDependencies(db)
	etiquetasDeps.RegisterRoutes(api)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	log.Printf("Swagger UI disponible en http://localhost:%s/swagger/index.html", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
