package main

import (
	"package/routes"
)

func main() {
	routes.InitializeRoutes()
	routes.StartServer()
}
