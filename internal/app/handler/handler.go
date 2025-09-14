package handler

import (
	"web_backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []repository.Device
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		devices, err = h.Repository.GetDevices()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		devices, err = h.Repository.GetDeviceByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	cartDevices, err := h.Repository.GetCart()
    cartCount := 0
    if err == nil {
        cartCount = len(cartDevices)
    }

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"devices": devices,
		"query": searchQuery,
		"cartCount": cartCount,
	})
}

func (h *Handler) GetDevice(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	device, err := h.Repository.GetDevice(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"device": device,
	})
}

func (h *Handler) GetCart(ctx *gin.Context) {
	var devices []repository.Device
	var err error

	devices, err = h.Repository.GetCart()
	if err != nil {
		logrus.Error(err)
		}

		
	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"service_devices": devices,
	})
}
