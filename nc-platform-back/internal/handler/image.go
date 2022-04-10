package handler

import (
	"demo-rest/internal/dto"
	"demo-rest/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	logger       *logrus.Logger
	imageService *service.ImageService
}

func NewImageHandler(service *service.ImageService) *ImageHandler {
	return &ImageHandler{
		logger:       logrus.New(),
		imageService: service,
	}
}

func (h *ImageHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.getUserImages)
	r.Post("/", h.upload)
	return r
}

func (h *ImageHandler) upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("img")
	if err != nil {
		h.logger.Error(err)
		return
	}
	defer file.Close()

	err = h.imageService.Upload(r.Context(), file, fileHeader)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ImageHandler) getUserImages(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		h.logger.Error(err)
	}

	pageSize, err := strconv.Atoi(r.FormValue("pageSize"))
	if err != nil {
		h.logger.Error(err)
	}

	pagingParam := dto.PagingParam{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	images, err := h.imageService.GetUserImages(r.Context(), pagingParam)
	if err != nil {
		h.logger.Error(err)
	}

	err = json.NewEncoder(w).Encode(images)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
