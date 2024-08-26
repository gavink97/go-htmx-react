package handlers

import (
	"github.com/gavink97/gavin-site/internal/store"
    "github.com/gavink97/gavin-site/internal/views"
	"net/http"
)

type PostRegisterHandler struct {
	userStore store.UserStore
}

type PostRegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewPostRegisterHandler(params PostRegisterHandlerParams) *PostRegisterHandler {
	return &PostRegisterHandler{
		userStore: params.UserStore,
	}
}

// dont allow users to create an account if one exists
func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(email, password)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		c := views.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := views.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
