package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"nc-platform-back/internal/service"
	"net/http"
)

type AuthHandler struct {
	logger      logrus.FieldLogger
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logrus.New(),
		authService: authService,
	}
}

type jwtResponse struct {
	Jwt string `json:"jwt"`
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.login)
	return r
}

func (h *AuthHandler) login(writer http.ResponseWriter, request *http.Request) {
	username, password, ok := request.BasicAuth()
	h.logger.Info(username, password, ok)
	if !ok {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := h.authService.Login(username, password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(writer).Encode(jwtResponse{Jwt: jwt})
}
