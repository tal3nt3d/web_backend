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

// GetDevices godoc
// @Summary Получить список устройств
// @Description Возвращает все устройства или фильтрует по названию
// @Tags devices
// @Produce json
// @Param device_title query string false "Название устройства для поиска"
// @Success 200 {array} serializer.DeviceJSON "Список устройств"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /devices [get]
func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []ds.Device
	var err error

	searchQuery := ctx.Query("device_title")
	if searchQuery == "" {
		devices, err = h.Repository.GetDevices()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		devices, err = h.Repository.GetDevicesByTitle(searchQuery)
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

// GetDevice godoc
// @Summary Получить устройство по ID
// @Description Возвращает информацию об устройстве по её идентификатору
// @Tags devices
// @Produce json
// @Param id path int true "ID устройства"
// @Success 200 {object} serializer.DeviceJSON "Данные устройства"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Устройство не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /device/{id} [get]
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

// CreateDevice godoc
// @Summary Создать новое устройство
// @Description Создает новое устройство и возвращает его данные
// @Tags devices
// @Accept json
// @Produce json
// @Param device body serializer.DeviceJSON true "Данные нового устройства"
// @Success 201 {object} serializer.DeviceJSON "Созданное устройство"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /device/create-device [post]
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

// DeleteDevice godoc
// @Summary Удалить устройство
// @Description Выполняет логическое удаление устройство по ID
// @Tags devices
// @Produce json
// @Param id path int true "ID устройства"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Устройство не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /device/{id}/delete-device [delete]
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

// EditDevice godoc
// @Summary Изменить данные устройства
// @Description Обновляет информацию об устройстве по ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "ID устройства"
// @Param device body serializer.DeviceJSON true "Новые данные устройства"
// @Success 200 {object} serializer.DeviceJSON "Обновленное устройство"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Устройство не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /device/{id}/edit-device [put]
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

// AddToAmperageApplication godoc
// @Summary Добавить устройство в расчёт
// @Description Добавляет устройство в заявку-черновик пользователя
// @Tags devices
// @Produce json
// @Param id path int true "ID устройства"
// @Success 200 {object} serializer.AmperageApplicationJSON "Расчёт с добавленным устройством"
// @Success 201 {object} serializer.AmperageApplicationJSON "Создан новый расчёт"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Устройство не найдено"
// @Failure 409 {object} map[string]string "Устройство уже в расчёте"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /device/{id}/add-to-amperage_application [post]
func (h *Handler) AddToAmperageApplication(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	amperage_application, created, err := h.Repository.GetAmperageApplicationDraft(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	amperage_application_id := amperage_application.Amperage_Application_ID

	device_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddToAmperageApplication(int(amperage_application_id), device_id)
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
		ctx.Header("Location", fmt.Sprintf("/amperage_application/%v", amperage_application.Amperage_Application_ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(amperage_application)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, serializer.AmperageApplicationToJSON(amperage_application, creatorLogin, moderatorLogin))
}

// AddPhoto godoc
// @Summary Загрузить изображение устройства
// @Description Загружает изображение для устройства и возвращает обновленные данные
// @Tags devices
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID устройства"
// @Param image formData file true "Изображение устройства"
// @Success 200 {object} map[string]interface{} "Статус загрузки и данные устройства"
// @Failure 400 {object} map[string]string "Неверный запрос или файл"
// @Failure 404 {object} map[string]string "Устройство не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /device/{id}/add-photo [post]
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
