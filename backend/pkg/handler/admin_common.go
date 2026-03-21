package handler

import (
	"errors"
	"net/http"

	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func getAdminID(c *gin.Context) (int64, bool) {
	role, _ := getRole(c)
	if role != "admin" {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return 0, false
	}

	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return 0, false
	}
	return userID, true
}

func isAdmin(c *gin.Context) bool {
	role, _ := getRole(c)
	return role == "admin"
}

func writeAdminApplicationErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrAppNotFound):
		writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
	case errors.Is(err, repository.ErrOperatorNotFound):
		writeErr(c, http.StatusNotFound, "not_found", "оператор не найден")
	case errors.Is(err, repository.ErrInvalidStatusCode):
		writeErr(c, http.StatusBadRequest, "validation_error", "неверный status_code")
	case errors.Is(err, repository.ErrFinalApplication):
		writeErr(c, http.StatusBadRequest, "validation_error", "нельзя изменить финальную заявку")
	default:
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
	}
}

func writeAdminUserErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		writeErr(c, http.StatusNotFound, "not_found", "пользователь не найден")
	case errors.Is(err, repository.ErrRoleNotFound):
		writeErr(c, http.StatusBadRequest, "validation_error", "неверный role_code")
	case errors.Is(err, repository.ErrSelfActionForbidden):
		writeErr(c, http.StatusForbidden, "forbidden", "нельзя выполнить это действие над своим аккаунтом")
	case errors.Is(err, repository.ErrUserHasRelations):
		writeErr(c, http.StatusConflict, "conflict", "пользователь связан с существующими данными")
	default:
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
	}
}

func writeAdminStatusErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrStatusNotFound):
		writeErr(c, http.StatusNotFound, "not_found", "статус не найден")
	case errors.Is(err, repository.ErrStatusExists):
		writeErr(c, http.StatusConflict, "conflict", "статус с таким code уже существует")
	case errors.Is(err, repository.ErrStatusProtected):
		writeErr(c, http.StatusForbidden, "forbidden", "системный статус нельзя изменять или удалять")
	case errors.Is(err, repository.ErrStatusInUse):
		writeErr(c, http.StatusConflict, "conflict", "статус используется в заявках")
	default:
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
	}
}
