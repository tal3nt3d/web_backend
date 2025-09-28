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
		"application_count": h.Repository.GetApplicationCount(),
		"Application_ID": h.Repository.GetActiveApplicationID(),
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

func (h *Handler) GetApplication(ctx *gin.Context) {
   	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		logrus.Error(err)
	}

	isDraft, err := h.Repository.IsDraftApplication(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
		return
	}

    applicationItems, err := h.Repository.GetApplication(id)
    if err != nil {
        logrus.Error(err)
	}

    ctx.HTML(http.StatusOK, "cart.html", gin.H{
        "application": applicationItems,
		"Application_ID": id,
    })
}

func (h *Handler) AddToApplication(ctx *gin.Context) {
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

func (h *Handler) DeleteApplication(ctx *gin.Context) {
	appIDStr := ctx.PostForm("application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteApplication(uint(appID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
