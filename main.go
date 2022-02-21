package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
	HandlerLayer "gofr-crud/handlers/users"
	serviceLayer "gofr-crud/services/users"
	"gofr-crud/stores/users"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false
	store := users.New()
	service := serviceLayer.New(store)
	ht := HandlerLayer.New(service)

	app.GET("/user/{id}", ht.GetByID)
	app.GET("/users", ht.GetAll)
	app.POST("/user", ht.Create)
	app.PUT("/user", ht.Update)
	app.DELETE("/user/{id}", ht.Delete)
	app.GET("/json", handler.JSONHandler)
	app.Server.HTTP.Port = 5000

	fmt.Println("Listening to PORT 5000")
	app.Start()
}
