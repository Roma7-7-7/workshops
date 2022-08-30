package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// ItemsList default implementation.
func ItemsList(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("items/list.html"))
}

// ItemsIndex default implementation.
func ItemsIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("items/index.html"))
}
