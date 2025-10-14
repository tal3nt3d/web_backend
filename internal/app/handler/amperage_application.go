package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAllAmperageApplications godoc
// @Summary Получить список заявок на расчёт
// @Description Возвращает заявки с возможностью фильтрации по датам и статусу
// @Tags amperage_applications
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус заявки"
// @Success 200 {array} serializer.AmperageApplicationJSON "Список заявок"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/all-amperage_applications [get]
func (h *Handler) GetAllAmperageApplications(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}
	fmt.Println(fromDate)

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	amperage_applications, err := h.Repository.GetAllAmperageApplications(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	amperage_applications = h.filterAuthorizedAmperageApplications(amperage_applications, ctx)
	resp := make([]serializer.AmperageApplicationJSON, 0, len(amperage_applications))
	for _, c := range amperage_applications {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, serializer.AmperageApplicationToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetAmperageApplicationCart godoc
// @Summary Получить корзину расчёта
// @Description Возвращает информацию о текущей заявке-черновике на расчёт пользователя
// @Tags amperage_applications
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины заявки-черновика"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/amperage_application-cart [get]
func (h *Handler) GetAmperageApplicationCart(ctx *gin.Context){
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	devices_count := h.Repository.GetAmperageApplicationCount(userID)

	if devices_count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"devices_count": devices_count,
		})
		return
	}

	amperage_application, err := h.Repository.CheckCurrentAmperageApplicationDraft(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":          "no_draft",
				"devices_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          amperage_application.Amperage_Application_ID,
		"devices_count": h.Repository.GetAmperageApplicationCount(amperage_application.Creator_ID),
	})
}

// GetAmperageApplication godoc
// @Summary Получить заявку по ID
// @Description Возвращает полную информацию о заявке
// @Tags amperage_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{} "Данные заявки с устройствами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/{id} [get]
func (h *Handler) GetAmperageApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	devices, amperage_application, err := h.Repository.GetAmperageApplicationDevices(id)
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

	resp := make([]serializer.DeviceJSON, 0, len(devices))
	for _, r := range devices {
		resp = append(resp, serializer.DeviceToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(amperage_application)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	amperage_applicationDevices, _ := h.Repository.GetDevicesAmperageApplications(int(amperage_application.Amperage_Application_ID))

	resp2 := make([]serializer.AmperageApplicationDeviceJSON, 0, len(amperage_applicationDevices))
	for _, r := range amperage_applicationDevices{
		resp2 = append(resp2, serializer.AmperageApplicationDeviceToJSON(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"amperage_application": serializer.AmperageApplicationToJSON(amperage_application, creatorLogin, moderatorLogin),
		"devices":   resp,
		"amperage_applicationDevices": resp2,
	})
}

// FormAmperageApplication godoc
// @Summary Сформировать заявку
// @Description Переводит заявку в статус "formed"
// @Tags amperage_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} serializer.AmperageApplicationJSON "Сформированная заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/{id}/form-amperage_application [put]
func (h *Handler) FormAmperageApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "formed"

	amperage_application, err := h.Repository.FormAmperageApplication(id, status)
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

// EditAmperageApplication godoc
// @Summary Изменить заявку
// @Description Обновляет данные заявки
// @Tags amperage_applications
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param amperage_applicatio body serializer.AmperageApplicationJSON true "Новые данные заявки"
// @Success 200 {object} serializer.AmperageApplicationJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/{id}/edit-amperage_application [put]
func (h *Handler) EditAmperageApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var amperage_applicationJSON serializer.AmperageApplicationJSON
	if err := ctx.BindJSON(&amperage_applicationJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	amperage_application, err := h.Repository.EditAmperageApplication(id, amperage_applicationJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
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

// DeleteAmperageApplication godoc
// @Summary Удалить заявку
// @Description Выполняет логическое удаление заявки
// @Tags amperage_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/{id}/delete-amperage_application [delete]
func (h *Handler) DeleteAmperageApplication(ctx *gin.Context){
	idStr := ctx.Param("id")
	amperage_application_id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "deleted"
	
	_, err = h.Repository.FormAmperageApplication(amperage_application_id, status)
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Amperage application deleted"})
}

// FinishAmperageApplication godoc
// @Summary Завершить заявку
// @Description Изменяет статус заявки (только для модераторов)
// @Tags amperage_applications
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param status body serializer.StatusJSON true "Новый статус"
// @Success 200 {object} serializer.AmperageApplicationJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /amperage_application/{id}/finish-amperage_application [put]
func (h *Handler) FinishAmperageApplication(ctx *gin.Context) {
	userID, err := getUserID(ctx)
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            h.errorHandler(ctx, http.StatusNotFound, err)
        } else {
            h.errorHandler(ctx, http.StatusInternalServerError, err)
        }
        return
    }
    
    if !user.IsModerator {
        h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
        return
    }

	amperage_application, err := h.Repository.FinishAmperageApplication(id, statusJSON.Status, userID)
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

func (h *Handler) filterAuthorizedAmperageApplications(amperage_applicaions []ds.AmperageApplication, ctx *gin.Context) []ds.AmperageApplication {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.AmperageApplication{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return []ds.AmperageApplication{}
	}
	if err != nil {
		return []ds.AmperageApplication{}
	}

	if user.IsModerator {
		return amperage_applicaions
	}

	var userAmperage_Applications []ds.AmperageApplication
    for _, amperage_application := range amperage_applicaions {
        fmt.Println(amperage_application.Amperage_Application_ID)
        if amperage_application.Creator_ID == userID {
            userAmperage_Applications = append(userAmperage_Applications, amperage_application)
        }
    }
    
    return userAmperage_Applications

}

func (h *Handler) hasAccessToAmperageApplication(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}