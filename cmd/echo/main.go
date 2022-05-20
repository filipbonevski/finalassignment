package main

import (
	"final/cmd"
	"final/cmd/echo/handlers"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"
)

//start the app with go run cmd/echo/main.go cmd/echo/db.go

func main() {
	db := initDB("data.db")
	migrate(db)
	CreateUser(db, "filipb", "blabla")
	CreateUser(db, "dada", "blabla")

	router := echo.New()

	auth := router.Group("/api")
	auth.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		a, err := handlers.AuthenticateUser(db, username, password)
		return a, err
	}))

	auth.GET("/lists/:id/tasks", handlers.GetTasks(db))
	auth.POST("/lists/:id/tasks", handlers.CreateTask(db))
	auth.PATCH("/tasks/:id", handlers.UpdateTask(db))
	auth.DELETE("/tasks/:id", handlers.DeleteTask(db))

	auth.GET("/lists", handlers.GetLists(db))
	auth.POST("/lists", handlers.CreateList(db))
	auth.DELETE("/lists/:id", handlers.DeleteList(db))

	auth.GET("/list/export", handlers.ExportTasks(db))
	auth.GET("/weather", handlers.GetWeather())

	// Do not touch this line!
	log.Fatal(http.ListenAndServe(":3000", cmd.CreateCommonMux(router)))
}
