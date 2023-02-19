package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber"
)

func writeError(c *fiber.Ctx, err error) error {
	{
		var httpErr *httpResponseMessage
		if errors.As(err, &httpErr) {
			return c.Status(httpErr.Status()).JSON(
				&struct {
					Status  int
					Message string
				}{
					httpErr.Status(),
					httpErr.Message(),
				})

		}
	}
	return c.Status(http.StatusInternalServerError).JSON(
		&struct {
			Status  int
			Message string
		}{
			http.StatusInternalServerError,
			err.Error(),
		})

}

func writeSuccessJSON(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(data)
}
