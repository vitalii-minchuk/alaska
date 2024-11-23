package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/vitalii-minchuk/alaska/services/auth"
	"github.com/vitalii-minchuk/alaska/types"
	"github.com/vitalii-minchuk/alaska/utils"
)

type Handler struct {
	store      types.ProductStore
	orderStore types.OrderStore
	userStore  types.UserStore
}

func NewHandler(
	store types.ProductStore,
	orderStore types.OrderStore,
	userStore types.UserStore,
) *Handler {
	return &Handler{
		store:      store,
		orderStore: orderStore,
		userStore:  userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handleCheckout).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WrightError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}

	// get products
	products, err := h.store.GetProductsByID(productIds)
	if err != nil {
		utils.WrightError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WrightError(w, http.StatusBadRequest, err)
		return
	}

	utils.WrightJSON(w, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
