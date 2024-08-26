package handlers

import (
	"github.com/gavink97/gavin-site/internal/layouts"
    "github.com/gavink97/gavin-site/internal/views"
	"net/http"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := views.NotFound()
	err := layouts.Layout(c, "Not Found").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
