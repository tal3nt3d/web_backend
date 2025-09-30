package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/repository"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []ds.Device
	var err error

	searchQuery := ctx.Query("device_name")
	if searchQuery == "" {
		devices, err = h.Repository.GetDevices()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		devices, err = h.Repository.GetDevicesByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	resp := make([]serializer.DeviceJSON, 0, len(devices))
	for _, r := range devices {
		resp = append(resp, serializer.DeviceToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetDevice(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device, err := h.Repository.GetDevice(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.DeviceToJSON(*device))
}

func (h *Handler) CreateDevice(ctx *gin.Context) {
	var deviceJSON serializer.DeviceJSON
	if err := ctx.BindJSON(&deviceJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	device, err := h.Repository.CreateDevice(deviceJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/device/%v", device.Device_ID))
	ctx.JSON(http.StatusCreated, serializer.DeviceToJSON(device))
}

func (h *Handler) DeleteDevice(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteDevice(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *Handler) EditDevice(ctx *gin.Context) {
	var deviceJSON serializer.DeviceJSON
	if err := ctx.BindJSON(&deviceJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device, err := h.Repository.EditDevice(id, deviceJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.DeviceToJSON(device))
}

func (h *Handler) AddToApplication(ctx *gin.Context) {
	application, created, err := h.Repository.GetApplicationDraft(uint(h.Repository.GetUserID()))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	application_id := application.Application_ID

	device_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddToApplication(int(application_id), device_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	
	status := http.StatusOK
	
	if created {
		ctx.Header("Location", fmt.Sprintf("/application/%v", application.Application_ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(application)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, serializer.ApplicationToJSON(application, creatorLogin, moderatorLogin))
}

func (h *Handler) AddPhoto(ctx *gin.Context) {
	device_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	device, err := h.Repository.AddPhoto(ctx, device_id, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
		"device": serializer.DeviceToJSON(device),
	})
}
