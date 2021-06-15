package main

import (
	"blog/controllers"
	"blog/core"
	_ "blog/core"
	// _ "blog/models"
)

func main() {
	controllerLogger := core.Logger.With("module", "controller")
	blog := controllers.NewHTTPAPIServer(controllerLogger)
	blog.Run()
}
