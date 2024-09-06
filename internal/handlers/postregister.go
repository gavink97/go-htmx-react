package handlers

import (
	"net/http"

	"github.com/gavink97/gavin-site/internal/store"
	"github.com/gavink97/gavin-site/internal/views"
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

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := h.userStore.GetUser(email)
	if err != nil {
		err := h.userStore.CreateUser(email, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			c := views.RegisterError()
			if err := c.Render(r.Context(), w); err != nil {
				http.Error(w, "Failed to render content", http.StatusInternalServerError)
			}

			return
		}
	} else {
		c := views.AccountExistsError()
		err = c.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := views.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
