package handler

import (
	"net/http"
	"strconv"

	"practice_IstraNet/pkg/dto"
	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserApps(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return
	}

	resp, err := h.services.GetUserApplications(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUserApp(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return
	}

	id := c.Query("id")
	appID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || appID <= 0 {
		writeErr(c, http.StatusBadRequest, "validation_error", "invalid id")
		return
	}

	resp, err := h.services.GetUserApplicationByID(c.Request.Context(), userID, appID)
	if err != nil {
		writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteUserApp(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return
	}

	var req dto.DeleteApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	err := h.services.DeleteUserApplication(c.Request.Context(), userID, req.ID)
	if err != nil {
		writeErr(c, http.StatusForbidden, "forbidden", "нельзя удалить заявку")
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func (h *Handler) ChangeUserApp(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		writeErr(c, http.StatusUnauthorized, "unauthorized", "нет userID")
		return
	}

	var req dto.ChangeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	err := h.services.UpdateUserApplication(c.Request.Context(), userID, req)
	if err != nil {
		switch err {
		case repository.ErrAppNotFound:
			writeErr(c, http.StatusNotFound, "not_found", "заявка не найдена")
		case repository.ErrForbidden:
			writeErr(c, http.StatusForbidden, "forbidden", "нельзя изменить чужую заявку")
		default:
			writeErr(c, http.StatusInternalServerError, "internal_error", "ошибка изменения")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "updated"})
}

