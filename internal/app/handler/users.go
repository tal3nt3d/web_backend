package handler

import (
	"errors"
	"net/http"
	"fmt"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(ctx *gin.Context) {
	var userJSON serializer.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if userJSON.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if userJSON.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}

	user, err := h.Repository.CreateUser(userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Header("Location", fmt.Sprintf("/users/%v", user.User_ID))
	ctx.JSON(http.StatusCreated, serializer.UserToJSON(user))
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var userJSON serializer.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if userJSON.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if userJSON.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}

	user, err := h.Repository.SignIn(userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else if err.Error() == "wrong password" {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.UserToJSON(user))
}

func (h *Handler) GetInfo(ctx *gin.Context) {
	userID := h.Repository.GetUserID()
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
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
	ctx.JSON(http.StatusOK, serializer.UserToJSON(user))
}

func (h *Handler) EditInfo(ctx *gin.Context) {
	userID := h.Repository.GetUserID()
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	var userJSON serializer.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.EditInfo(userID, userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.UserToJSON(user))
}

func (h *Handler) SignOut(ctx *gin.Context) {
	h.Repository.SignOut()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}