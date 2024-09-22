package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/layouts"
	"github.com/gavink97/gavin-site/internal/middleware"
	"github.com/gavink97/gavin-site/internal/store"
	"github.com/gavink97/gavin-site/internal/views"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	subcookie := checkStatusCookie(r)
	user, ok := r.Context().Value(middleware.UserKey).(*store.User)

	if !ok {
		c := views.GuestIndex()

		err := layouts.Layout(c, "My website", subcookie).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		return
	}

	c := views.IndexAuth(user.Email)
	err := layouts.Layout(c, "My website", subcookie).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
