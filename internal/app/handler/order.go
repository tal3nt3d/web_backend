package handler

import (
	"net/http"
	"strconv"
	"web_backend/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []ds.Device
	var err error

	searchQuery := ctx.Query("query") 
	if searchQuery == "" {            
		devices, err = h.Repository.GetDevices()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		devices, err = h.Repository.GetDevicesByTitle(searchQuery) 
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"devices": devices,
		"query":  searchQuery, 
		"amperage_application_count": h.Repository.GetAmperageApplicationCount(),
		"Amperage_Application_ID": h.Repository.GetActiveAmperageApplicationID(),
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

func (h *Handler) GetAmperageApplication(ctx *gin.Context) {
   	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		logrus.Error(err)
	}

	isDraft, err := h.Repository.IsDraftAmperageApplication(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
		return
	}

    amperageApplicationItems, err := h.Repository.GetAmperageApplication(id)
    if err != nil {
        logrus.Error(err)
	}

    ctx.HTML(http.StatusOK, "cart.html", gin.H{
        "amperage_application": amperageApplicationItems,
		"Amperage_Application_ID": id,
    })
}

func (h *Handler) AddToAmperageApplication(ctx *gin.Context) {
    deviceIDStr := ctx.PostForm("device_id")
    deviceID, err := strconv.Atoi(deviceIDStr)
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

    creatorID := uint(1)

    err = h.Repository.AddDevice(uint(deviceID), creatorID)
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }

    ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteAmperageApplication(ctx *gin.Context) {
	appIDStr := ctx.PostForm("amperage_application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteAmperageApplication(uint(appID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
