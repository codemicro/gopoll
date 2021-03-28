package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) error {

	app.Get("/", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	// ------- NEW POLL -------

	app.Get("/new", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	app.Post("/new", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	// ------- VOTE IN POLL -------

	app.Get("/vote", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	app.Post("/vote/:pollId", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	// ------- VIEW POLL RESULTS -------

	app.Get("/results/:pollId", func(ctx *fiber.Ctx) error {
		// TODO
		return nil
	})

	return nil
}
