package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/repository"
	"errors"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/devices", h.GetDevices)
	router.GET("/api/device/:id", h.GetDevice)
	router.POST("/api/device/create-device", h.CreateDevice)
	router.PUT("/api/device/:id/edit-device", h.EditDevice)
	router.DELETE("/api/device/:id/delete-device", h.DeleteDevice)
	router.POST("/api/device/:id/add-to-amperage_application", h.AddToAmperageApplication)
	router.POST("/api/device/:id/add-photo", h.AddPhoto)

	router.GET("/api/amperage_application/amperage_application-cart", h.GetAmperageApplicationCart)
	router.GET("/api/amperage_application/all-amperage_applications", h.GetAllAmperageApplications)
	router.GET("/api/amperage_application/:id", h.GetAmperageApplication)
	router.PUT("/api/amperage_application/:id/edit-amperage_application", h.EditAmperageApplication)
	router.PUT("/api/amperage_application/:id/form-amperage_application", h.FormAmperageApplication)
	router.PUT("/api/amperage_application/:id/finish-amperage_application", h.FinishAmperageApplication)
	router.DELETE("/api/amperage_application/:id/delete-amperage_application", h.DeleteAmperageApplication) 

	router.DELETE("/api/dev_app/:device_id/:amperage_application_id", h.DeleteDeviceFromAmperageApplication)
	router.PUT("/api/dev_app/:device_id/:amperage_application_id", h.EditDeviceFromAmperageApplication)

	router.POST("/api/users/signup", h.CreateUser)
	router.GET("/api/users/info", h.GetInfo)
	router.PUT("/api/users/info", h.EditInfo)
	router.POST("/api/users/signin", h.SignIn)
	router.POST("/api/users/signout", h.SignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())

	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}

