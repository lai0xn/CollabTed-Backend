package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CollabTed/CollabTed-Backend/internal/auth/usecase"
	"github.com/CollabTed/CollabTed-Backend/pkg/dto"
	"github.com/CollabTed/CollabTed-Backend/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(useCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user dto.UserD
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := utils.Validate.Struct(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errors := utils.BuildError(err.(validator.ValidationErrors))
		if err := json.NewEncoder(w).Encode(errors); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := h.useCase.Register(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.useCase.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
