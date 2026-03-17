package handler

import (
	"errors"
	"net/http"
	"strconv"

	"practice_IstraNet/pkg/dto"
	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OperatorGetApps(c *gin.Context) {
	if !isOperator(c) {
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

func (h *Handler) OperatorGetApp(c *gin.Context) {
	if !isOperator(c) {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return
	}

	appID, ok := getPositiveIDQuery(c, "id")
	if !ok {
		return
	}

	resp, err := h.services.GetOperatorApplicationByID(c.Request.Context(), appID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
		} else {
			writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) OperatorTakeApp(c *gin.Context) {
	operatorID, ok := getOperatorID(c)
	if !ok {
		return
	}

	var req dto.TakeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.TakeApplication(c.Request.Context(), operatorID, req.ID); err != nil {
		writeOperatorErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "taken"})
}

func (h *Handler) OperatorChangeStatus(c *gin.Context) {
	operatorID, ok := getOperatorID(c)
	if !ok {
		return
	}

	var req dto.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.ChangeApplicationStatus(c.Request.Context(), operatorID, req); err != nil {
		writeOperatorErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "status_changed"})
}

func (h *Handler) OperatorGetHistory(c *gin.Context) {
	if !isOperator(c) {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return
	}

	appID, ok := getPositiveIDQuery(c, "id")
	if !ok {
		return
	}

	resp, err := h.services.GetApplicationHistory(c.Request.Context(), appID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
		} else {
			writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) OperatorCloseApp(c *gin.Context) {
	operatorID, ok := getOperatorID(c)
	if !ok {
		return
	}

	var req dto.CloseApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.CloseApplication(c.Request.Context(), operatorID, req); err != nil {
		writeOperatorErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "closed"})
}

func getOperatorID(c *gin.Context) (int64, bool) {
	role, _ := getRole(c)
	if role != "operator" {
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

func isOperator(c *gin.Context) bool {
	role, _ := getRole(c)
	return role == "operator"
}

func getPositiveIDQuery(c *gin.Context, key string) (int64, bool) {
	idStr := c.Query(key)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeErr(c, http.StatusBadRequest, "validation_error", "invalid id")
		return 0, false
	}
	return id, true
}

func writeOperatorErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrAppNotFound):
		writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
	case errors.Is(err, repository.ErrForbidden):
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа к заявке")
	case errors.Is(err, repository.ErrAlreadyAssigned):
		writeErr(c, http.StatusConflict, "conflict", "заявка уже взята в работу")
	case errors.Is(err, repository.ErrInvalidStatusTransition), errors.Is(err, repository.ErrInvalidStatusCode):
		writeErr(c, http.StatusBadRequest, "validation_error", "недопустимое изменение статуса")
	default:
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
	}
}
