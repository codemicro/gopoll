package main

import (
	"github.com/codemicro/gopoll/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	if err := routes.Setup(app); err != nil {
		panic(err)
	}

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
