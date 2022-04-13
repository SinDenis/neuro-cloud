package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"nc-platform-back/internal/domain"
	"nc-platform-back/internal/service"
	"net/http"
)

type RegisterHandler struct {
	logger          logrus.FieldLogger
	registerService *service.RegisterService
}

func NewRegisterHandler(registerService *service.RegisterService) *RegisterHandler {
	return &RegisterHandler{
		logger:          logrus.New(),
		registerService: registerService,
	}
}

func (h *RegisterHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.register)
	return r
}

func (h *RegisterHandler) register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error(err)
		return
	}

	fmt.Println(user)
	err = h.registerService.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
