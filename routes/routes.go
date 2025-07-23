package routes

import (
	"github.com/Antonio-Jacal/papeleria-backend.git/controllers"
	"github.com/Antonio-Jacal/papeleria-backend.git/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Rutas públicas
	api.POST("/login", controllers.Login)
	api.GET("/listas", controllers.GetList)

	// Ruta protegida solo para admin y developxz
	protected := api.Group("/protegida")
	protected.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "develop"))
	protected.GET("", func(c *gin.Context) {
		rol, _ := c.Get("rol")
		userId, _ := c.Get("userId")
		c.JSON(200, gin.H{"mensaje": "Bienvenido la ruta protegida", "rol": rol, "userId": userId})
	})

	registerList := api.Group("/registrolista")
	registerList.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "develop", "worker"))
	registerList.POST("", controllers.RegisterList)

	register := api.Group("register")
	register.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("develop"))
	register.POST("", controllers.Register)

	getListWithFilters := api.Group("/getlist")
	getListWithFilters.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "worker", "develop"))
	getListWithFilters.GET("", controllers.GetListWithFilters)

	getList := getListWithFilters.Group("/:id")
	getList.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "worker", "develop"))
	getList.GET("", controllers.GetPedidoID)

	reports := api.Group("/reportes")
	reports.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "develop", "worker"))

	reportCash := reports.Group("/resumen-facturacion")
	reportCash.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "develop"))
	reportCash.GET("", controllers.GetSummary)

	reportLabels := reports.Group("/resumen-etiquetas")
	reportLabels.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "worker", "develop"))
	reportLabels.GET("", controllers.GetSummaryLabels)

	reportFaltantes := reports.Group("/resumen-faltantes")
	reportFaltantes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "develop", "worker"))
	reportFaltantes.GET("", controllers.GetSummaryFaltantes)

	updateLabel := api.Group("/updateList")
	updateLabel.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "worker", "develop"))
	getLabelId := updateLabel.Group("/:id")
	getLabelId.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "worker", "develop"))
	getLabelId.PUT("", controllers.UpdateLabelResponse)
	getLabelId.POST("", controllers.UpdatePedido)

}
