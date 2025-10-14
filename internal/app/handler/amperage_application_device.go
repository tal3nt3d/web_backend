package handler

import (
	"errors"
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

// DeleteDeviceFromAmperageApplication godoc
// @Summary Удалить устройство из заявки
// @Description Удаляет связь устройства и заявки
// @Tags amperage_application_devices
// @Produce json
// @Param device_id path int true "ID устройства"
// @Param amperage_application_id path int true "ID заявки"
// @Success 200 {object} serializer.AmperageApplicationJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /dev_app/{device_id}/{amperage_application_id} [delete]
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

// EditDeviceFromAmperageApplication godoc
// @Summary Изменить данные устройства в заявке
// @Description Обновляет параметры устройства в конкретной заявке
// @Tags amperage_application_devices
// @Accept json
// @Produce json
// @Param device_id path int true "ID устройства"
// @Param amperage_application_id path int true "ID заявки"
// @Param data body serializer.AmperageApplicationDeviceJSON true "Новые данные"
// @Success 200 {object} serializer.AmperageApplicationDeviceJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /dev_app/{device_id}/{amperage_application_id} [put]
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