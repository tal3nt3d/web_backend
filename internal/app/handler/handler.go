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
	router.GET("/", h.GetDevices)
	router.GET("/device/:id", h.GetDevice)
	router.POST("/device/create-device", h.CreateDevice)
	router.PUT("/device/:id/edit-device", h.EditDevice)
	router.DELETE("/device/:id/delete-device", h.DeleteDevice)
	router.POST("/device/:id/add-to-application", h.AddToApplication)
	router.POST("/device/:id/add-photo", h.AddPhoto)

	router.GET("/application/application-cart", h.GetApplicationCart)
	router.GET("/application/all-applications", h.GetAllApplications)
	router.GET("/application/:id", h.GetApplication)
	router.PUT("/application/:id/edit-application", h.EditApplication)
	router.PUT("/application/:id/form-application", h.FormApplication)
	router.PUT("/application/:id/finish-application", h.FinishApplication)
	router.DELETE("/application/:id/delete-application", h.DeleteApplication) 

	router.DELETE("/dev_app/:device_id/:application_id/", h.DeleteDeviceFromApplication)
	router.PUT("/dev_app/:device_id/:application_id/", h.EditDeviceFromApplication)

	router.POST("/users/signup", h.CreateUser)
	router.GET("/users/info", h.GetInfo)
	router.PUT("/users/info", h.EditInfo)
	router.POST("/users/signin", h.SignIn)
	router.POST("/users/signout", h.SignOut)
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

