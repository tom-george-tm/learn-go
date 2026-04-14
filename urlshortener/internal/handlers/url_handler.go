package handlers

import (
	"errors"
	"net/http"

	"urlshortener/internal/service"

	"github.com/gofiber/fiber/v2"
)

// URLHandler handles HTTP requests for URLs.
type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/urls")
	api.Post("/", h.Shorten)
	api.Get("/", h.List)
	api.Get("/:code", h.GetStats)
	api.Delete("/:code", h.Delete)

	app.Get("/:code", h.Redirect)
}

func (h *URLHandler) Shorten(c *fiber.Ctx) error {
	var req service.ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	res, err := h.svc.Shorten(req)
	if err != nil {
		if errors.Is(err, service.ErrCodeTaken) {
			return errorResponse(c, http.StatusConflict, err.Error())
		}
		return errorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *URLHandler) Redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	url, err := h.svc.GetOriginal(code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrExpired) {
			return errorResponse(c, http.StatusGone, err.Error())
		}
		return errorResponse(c, http.StatusInternalServerError, "internal error")
	}

	return c.Redirect(url, http.StatusFound)
}

func (h *URLHandler) GetStats(c *fiber.Ctx) error {
	code := c.Params("code")
	res, err := h.svc.GetStats(code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, err.Error())
		}
		return errorResponse(c, http.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *URLHandler) Delete(c *fiber.Ctx) error {
	code := c.Params("code")
	if err := h.svc.Delete(code); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, err.Error())
		}
		return errorResponse(c, http.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "deleted",
	})
}

func (h *URLHandler) List(c *fiber.Ctx) error {
	urls, err := h.svc.List()
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(urls),
		"data":    urls,
	})
}

func errorResponse(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   msg,
	})
}
