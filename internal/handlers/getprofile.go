package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/layouts"
	"github.com/gavink97/gavin-site/internal/middleware"
	"github.com/gavink97/gavin-site/internal/store"
	"github.com/gavink97/gavin-site/internal/views"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(middleware.UserKey).(*store.User)

	if !ok {
		// http.Redirect()
		return
	}

	subcookie := checkStatusCookie(r)
	c := views.Profile(user.Email, subcookie)

	err := layouts.Layout(c, "User Profile", subcookie).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
