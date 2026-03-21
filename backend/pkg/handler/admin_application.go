package handler

import (
	"errors"
	"net/http"

	"practice_IstraNet/pkg/dto"
	"practice_IstraNet/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminGetApps(c *gin.Context) {
	if !isAdmin(c) {
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

func (h *Handler) AdminGetApp(c *gin.Context) {
	if !isAdmin(c) {
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

func (h *Handler) AdminAssignApp(c *gin.Context) {
	adminID, ok := getAdminID(c)
	if !ok {
		return
	}

	var req dto.AdminAssignApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.AssignApplication(c.Request.Context(), adminID, req); err != nil {
		writeAdminApplicationErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "assigned"})
}

func (h *Handler) AdminChangeStatus(c *gin.Context) {
	adminID, ok := getAdminID(c)
	if !ok {
		return
	}

	var req dto.AdminChangeApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.ChangeApplicationStatusByAdmin(c.Request.Context(), adminID, req); err != nil {
		writeAdminApplicationErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "status_changed"})
}

func (h *Handler) AdminGetHistory(c *gin.Context) {
	if !isAdmin(c) {
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

func (h *Handler) AdminDeleteApp(c *gin.Context) {
	if _, ok := getAdminID(c); !ok {
		return
	}

	var req dto.AdminDeleteApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.DeleteApplicationByAdmin(c.Request.Context(), req.ID); err != nil {
		writeAdminApplicationErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}
