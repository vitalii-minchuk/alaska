package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/vitalii-minchuk/alaska/config"
	"github.com/vitalii-minchuk/alaska/services/auth"
	"github.com/vitalii-minchuk/alaska/types"
	"github.com/vitalii-minchuk/alaska/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WrightJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("user with emil %s already exists", payload.Email))
		return
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  hashedPassword,
		Email:     payload.Email,
	})
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WrightJSON(w, http.StatusCreated, nil)
}
