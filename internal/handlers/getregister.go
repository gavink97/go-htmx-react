package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/layouts"
	"github.com/gavink97/gavin-site/internal/views"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler() *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := views.RegisterPage()
	subcookie := checkStatusCookie(r)

	err := layouts.Layout(c, "My website", subcookie).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}
