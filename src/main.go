package main

import "application"

func main() {
	app := application.App{}
	app.Initialize(
		"mateo",
		"prueba")
	app.Run("8082")
}
