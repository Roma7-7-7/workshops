package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// CategoriesList default implementation.
func CategoriesList(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("categories/list.html"))
}

// CategoriesIndex default implementation.
func CategoriesIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("categories/index.html"))
}
