package handler

import (
	"net/http"

	"practice_IstraNet/pkg/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminGetUsers(c *gin.Context) {
	if !isAdmin(c) {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return
	}

	resp, err := h.services.GetUsers(c.Request.Context())
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "internal_error", "внутренняя ошибка")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AdminGetUser(c *gin.Context) {
	if !isAdmin(c) {
		writeErr(c, http.StatusForbidden, "forbidden", "нет доступа")
		return
	}

	userID, ok := getPositiveIDQuery(c, "id")
	if !ok {
		return
	}

	resp, err := h.services.GetUserByIDForAdmin(c.Request.Context(), userID)
	if err != nil {
		writeAdminUserErr(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AdminChangeRole(c *gin.Context) {
	adminID, ok := getAdminID(c)
	if !ok {
		return
	}

	var req dto.AdminChangeUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.ChangeUserRole(c.Request.Context(), adminID, req); err != nil {
		writeAdminUserErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "role_changed"})
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	adminID, ok := getAdminID(c)
	if !ok {
		return
	}

	var req dto.AdminDeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.DeleteUserByAdmin(c.Request.Context(), adminID, req.UserID); err != nil {
		writeAdminUserErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}
