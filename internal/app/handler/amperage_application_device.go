package handler

import (
	"errors"
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteDeviceFromAmperageApplication(ctx *gin.Context) {
	amperage_application_id, err := strconv.Atoi(ctx.Param("amperage_application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	amperage_application, err := h.Repository.DeleteDeviceFromAmperageApplication(amperage_application_id, device_id)
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

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(amperage_application)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.AmperageApplicationToJSON(amperage_application, creatorLogin, moderatorLogin))
}

func (h *Handler) EditDeviceFromAmperageApplication(ctx *gin.Context) {
	amperage_application_id, err := strconv.Atoi(ctx.Param("amperage_application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var amperage_applicationDeviceJSON serializer.AmperageApplicationDeviceJSON
	if err := ctx.BindJSON(&amperage_applicationDeviceJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	amperage_applicationDevice, err := h.Repository.EditDeviceFromAmperageApplication(amperage_application_id, device_id, amperage_applicationDeviceJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.AmperageApplicationDeviceToJSON(amperage_applicationDevice))
}