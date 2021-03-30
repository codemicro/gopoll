package routes

import (
	"github.com/codemicro/gopoll/internal/pages"
	"github.com/codemicro/gopoll/internal/webRes"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) error {

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Type("html").SendString(pages.Homepage())
	})

	// ------- NEW POLL -------

	app.Get("/new", func(ctx *fiber.Ctx) error {
		return ctx.Type("html").SendString(pages.NewPoll())
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

	// ------- RESOURCES -------

	// TODO: compression
	app.Get("/res/main.css", func(ctx *fiber.Ctx) error {
		return ctx.Type("css").Send(webRes.MustAsset("main.css"))
	})

	return nil
}
