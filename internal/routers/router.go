package routers

import (
	"stock-flow/internal/controllers"
	"stock-flow/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Controllers
	authCtrl := new(controllers.AuthController)
	matCtrl := new(controllers.MaterialController)
	invCtrl := controllers.NewInventoryController()
	outCtrl := new(controllers.OutboundController)

	// Public
	auth := r.Group("/auth")
	{
		auth.POST("/login", authCtrl.Login)
		auth.POST("/register", authCtrl.Register)
	}

	// Protected
	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth())
	{
		// Material (Admin/Keeper)
		mat := api.Group("/materials")
		mat.Use(middleware.RoleAuth("Admin", "Keeper"))
		{
			mat.POST("", matCtrl.Create)
			mat.GET("", matCtrl.List)
			mat.DELETE("/:id", matCtrl.Delete)
		}

		// Inventory
		inv := api.Group("/inventory")
		{
			// Inbound (Keeper)
			inv.POST("/inbound", middleware.RoleAuth("Admin", "Keeper"), invCtrl.Inbound)
			inv.POST("/import", middleware.RoleAuth("Admin", "Keeper"), invCtrl.BatchImport)
			inv.DELETE("/:id", middleware.RoleAuth("Admin", "Keeper"), invCtrl.Delete)

			// List (All)
			inv.GET("", invCtrl.List)

			// Recommend (All)
			inv.GET("/recommend", invCtrl.RecommendedBatches)
		}

		// Outbound
		out := api.Group("/outbound")
		{
			out.POST("/apply", outCtrl.Apply)
			out.GET("/my", outCtrl.List)
			out.PUT("/:id/status", outCtrl.UpdateStatus)
		}
	}

	return r
}
