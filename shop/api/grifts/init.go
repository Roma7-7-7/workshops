package grifts

import (
	"github.com/Roma7-7-7/workshops/shop/api/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
