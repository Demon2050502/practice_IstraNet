package handler

import (
	"net/http"

	"practice_IstraNet/pkg/dto"
	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateApplication(c *gin.Context) {
	role, _ := getRole(c)
	if role != "user" {
		writeErr(c, http.StatusForbidden, "forbidden", "заявку может создать только пользователь")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return
	}

	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	resp, err := h.services.CreateApplication(c.Request.Context(), userID, req)
	if err != nil {
		switch err {
		case repository.ErrPriorityNotFound:
			writeErr(c, http.StatusBadRequest, "validation_error", "неверный priority_code")
		case repository.ErrCategoryNotFound:
			writeErr(c, http.StatusBadRequest, "validation_error", "category_id не найден")
		case repository.ErrStatusNotFound:
			writeErr(c, http.StatusInternalServerError, "internal_error", "в БД нет статуса 'new'")
		default:
			writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAllApplications(c *gin.Context) {
	role, _ := getRole(c)
	if role != "operator" && role != "admin" {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return
	}

	resp, err := h.services.GetAllApplications(c.Request.Context())
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		return
	}

	c.JSON(http.StatusOK, resp)
}
