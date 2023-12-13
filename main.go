package main

import (
	"github.com/amannlalwani/Notes-app-using-gofr/handler"
	"github.com/amannlalwani/Notes-app-using-gofr/store"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	s := store.New()
	h := handler.New(s)

	app.Server.ValidateHeaders = false

	// specifying the different routes supported by this service
	app.GET("/notes", h.Get)
	app.POST("/notes", h.Create)
	app.PUT("/notes/{id}", h.Update)
	app.DELETE("/notes/{id}", h.Delete)

	app.Server.HTTP.Port = 3000
	app.Server.MetricsPort = 2113

	app.Start()
}
