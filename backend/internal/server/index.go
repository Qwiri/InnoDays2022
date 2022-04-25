package server

import "github.com/gofiber/fiber/v2"

func (s *Server) routeIndex(c *fiber.Ctx) (err error) {
	return c.SendString("Hello mom")
}
