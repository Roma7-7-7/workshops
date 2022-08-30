package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// OrdersCreate default implementation.
func OrdersCreate(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("orders/create.html"))
}

// OrdersUpdate default implementation.
func OrdersUpdate(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("orders/update.html"))
}
