package handler

import (
	"net/http"

	"practice_IstraNet/pkg/dto"
	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
    var req dto.SignUpRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
        return
    }

    resp, err := h.services.SignUp(c.Request.Context(), req)
    if err != nil {
        switch err {
        case repository.ErrUserExists:
            writeErr(c, http.StatusConflict, "conflict", "пользователь с таким email уже существует")
        case repository.ErrRoleNotFound:
            writeErr(c, http.StatusBadRequest, "validation_error", "неизвестная роль")
        default:
            writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
        }
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) SignIn(c *gin.Context) {
    var req dto.SignInRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
        return
    }

    resp, err := h.services.SignIn(c.Request.Context(), req)
    if err != nil {
        switch err {
        case dto.ErrInvalidCredentials, repository.ErrUserNotFound:
            writeErr(c, http.StatusUnauthorized, "unauthorized", "неверный email или пароль")
        case repository.ErrUserInactive:
            writeErr(c, http.StatusForbidden, "forbidden", "пользователь отключён")
        default:
            writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
        }
        return
    }

    c.JSON(http.StatusOK, resp)
}

func writeErr(c *gin.Context, status int, errType, msg string) {
    c.JSON(status, dto.ErrorResponse{
        Error:   errType,
        Message: msg,
    })
}
