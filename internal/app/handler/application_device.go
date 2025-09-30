package handler

import (
	"errors"
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteDeviceFromApplication(ctx *gin.Context) {
	application_id, err := strconv.Atoi(ctx.Param("application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	application, err := h.Repository.DeleteDeviceFromApplication(application_id, device_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(application)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.ApplicationToJSON(application, creatorLogin, moderatorLogin))
}

func (h *Handler) EditDeviceFromApplication(ctx *gin.Context) {
	application_id, err := strconv.Atoi(ctx.Param("application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var applicationDeviceJSON serializer.ApplicationDeviceJSON
	if err := ctx.BindJSON(&applicationDeviceJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	applicationDevice, err := h.Repository.EditDeviceFromApplication(application_id, device_id, applicationDeviceJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.ApplicationDeviceToJSON(applicationDevice))
}