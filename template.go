package ace

import (
	"net/http"
)

type Context map[string]interface{}

// Renderer is the interface that lets you use any HTML renderer
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}

// HtmlTemplate sets the renderer to use for HTML templates
func (a *Ace) HtmlTemplate(render Renderer) {
	a.render = render
}
