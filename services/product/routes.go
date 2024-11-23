package product

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/vitalii-minchuk/alaska/types"
	"github.com/vitalii-minchuk/alaska/utils"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.createProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.getProducts).Methods(http.MethodGet)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	err := h.store.CreateProduct(&types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Quantity:    payload.Quantity,
		Image:       payload.Image,
		Price:       payload.Price,
	})
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WrightJSON(w, http.StatusCreated, nil)
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WrightJSON(w, http.StatusOK, ps)
}
