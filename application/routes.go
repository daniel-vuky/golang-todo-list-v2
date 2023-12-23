package application

import (
	"github.com/daniel-vuky/golang-todo-list-and-chat/handler"
	"github.com/daniel-vuky/golang-todo-list-and-chat/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoadRoutes load all the routes of application
func (app *App) LoadRoutes() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	LoadIndexTemplate(router)
	LoadItemRoutes(app, router)

	app.router = router
}

// LoadIndexTemplate load the index template
func LoadIndexTemplate(router *gin.Engine) {
	// Set the HTML templates directory
	router.Static("/static", "./public/static")
	router.LoadHTMLGlob("./public/templates/*")

	// Define a route to render the index template
	router.GET("/", func(c *gin.Context) {
		// Render the HTML template with data
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}

// LoadItemRoutes load all the items api routes
func LoadItemRoutes(app *App, router *gin.Engine) {
	itemsHandler := &handler.Items{
		Repository: &repository.ItemsRepository{
			Db: app.rdb,
		},
	}
	itemGroup := router.Group("/items")
	{
		itemGroup.GET("/", itemsHandler.List)
		itemGroup.POST("/", itemsHandler.Create)
		itemGroup.GET("/:id", itemsHandler.GetByID)
		itemGroup.PUT("/:id", itemsHandler.UpdateByID)
		itemGroup.DELETE("/:id", itemsHandler.DeleteByID)
	}
}
