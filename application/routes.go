package application

import (
	"github.com/daniel-vuky/golang-todo-list-and-chat/handler"
	"github.com/daniel-vuky/golang-todo-list-and-chat/repository"
	"github.com/gin-gonic/gin"
	ginSession "github.com/go-session/gin-session"
	"net/http"
)

// LoadRoutes load all the routes of application
func (app *App) LoadRoutes() {
	router := gin.Default()
	router.Use(ginSession.New())

	// Set the HTML templates directory
	router.Static("/static", "./public/static")
	router.LoadHTMLGlob("./public/templates/*")

	usersHandler := &handler.Users{
		Repository: &repository.UsersRepository{
			Db: app.rdb,
		},
		Auth: &repository.AuthRepository{
			Db: app.rdb,
		},
	}

	LoadAuthRoutes(app, router, usersHandler)
	LoadItemRoutes(app, router, usersHandler)

	app.router = router
}

// LoadAuthRoutes load the auth template
func LoadAuthRoutes(app *App, router *gin.Engine, usersHandler *handler.Users) {

	router.GET("/", usersHandler.AuthMiddleware, func(c *gin.Context) {
		username := usersHandler.GetUsernameFromContext(c)
		c.HTML(http.StatusOK, "index.html", gin.H{"username": username})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	router.POST("/login", usersHandler.Login)
	router.GET("/logout", usersHandler.Logout)
	router.POST("/register", usersHandler.Register)
}

// LoadItemRoutes load all the items api routes
func LoadItemRoutes(app *App, router *gin.Engine, usersHandler *handler.Users) {
	itemsHandler := &handler.Items{
		Repository: &repository.ItemsRepository{
			Db: app.rdb,
		},
	}
	itemGroup := router.Group("/items")
	{
		itemGroup.GET("/", usersHandler.AuthMiddleware, itemsHandler.List)
		itemGroup.POST("/", usersHandler.AuthMiddleware, itemsHandler.Create)
		itemGroup.GET("/:id", usersHandler.AuthMiddleware, itemsHandler.GetByID)
		itemGroup.PUT("/:id", usersHandler.AuthMiddleware, itemsHandler.UpdateByID)
		itemGroup.DELETE("/:id", usersHandler.AuthMiddleware, itemsHandler.DeleteByID)
	}
}
