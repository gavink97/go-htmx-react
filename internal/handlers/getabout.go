package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/layouts"
	"github.com/gavink97/gavin-site/internal/views"
)

type AboutHandLer struct {
}

func NewAboutHandler() *AboutHandLer {
	return &AboutHandLer{}
}

func (h *AboutHandLer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := views.About()

	subcookie := checkStatusCookie(r)
	err := layouts.Layout(c, "My website", subcookie).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
