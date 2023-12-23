package application

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	router *gin.Engine
	rdb    *sql.DB
}

// New create new application
// Launch database connection and load the routes
// Return App
func New() *App {
	sourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER_NAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE_NAME"),
	)
	db, connectedErr := sql.Open("mysql", sourceName)
	if connectedErr != nil {
		log.Fatalf(fmt.Sprintf("Can not connect to mysql server, %s", connectedErr.Error()))
	}

	app := &App{
		rdb: db,
	}

	app.LoadRoutes()

	return app
}

// Start method of App
// Starting server
// Add listener the cancel event
// Return nil or error
func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APPLICATION_PORT")),
		Handler: app.router,
	}

	pingErr := app.rdb.Ping()
	if pingErr != nil {
		return fmt.Errorf("Can not ping to database, %s", pingErr.Error())
	}
	defer func() {
		if err := app.rdb.Close(); err != nil {
			fmt.Printf("Fail to close connect to data base, %s", err.Error())
		}
	}()

	fmt.Println("Starting server")

	serverError := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			serverError <- fmt.Errorf("Fail to start the server, %s", err.Error())
		}
	}()

	select {
	case err := <-serverError:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}
