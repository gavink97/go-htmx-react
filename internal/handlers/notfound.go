package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/layouts"
	"github.com/gavink97/gavin-site/internal/views"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := views.NotFound()
	subcookie := checkStatusCookie(r)

	err := layouts.Layout(c, "Not Found", subcookie).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
