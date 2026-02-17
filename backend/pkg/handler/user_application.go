package handler

import (
	"fmt"
	"practice_IstraNet/pkg/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserApps(c *gin.Context) {
	userID, _ := getUserID(c)
	resp, err := h.services.GetUserApplications(c.Request.Context(), userID)
	if err != nil {
		writeErr(c, 500, "internal_error", "внутренняя ошибка")
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) GetUserApp(c *gin.Context) {
	userID, _ := getUserID(c)
	id := c.Query("id")

	var appID int64
	fmt.Sscan(id, &appID)

	resp, err := h.services.GetUserApplicationByID(c.Request.Context(), userID, appID)
	if err != nil {
		writeErr(c, 404, "not_found", "заявка не найдена")
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) DeleteUserApp(c *gin.Context) {
	userID, _ := getUserID(c)

	var req dto.DeleteApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, 400, "validation_error", err.Error())
		return
	}

	err := h.services.DeleteUserApplication(c.Request.Context(), userID, req.ID)
	if err != nil {
		writeErr(c, 403, "forbidden", "нельзя удалить заявку")
		return
	}

	c.JSON(200, gin.H{"result": "deleted"})
}

func (h *Handler) ChangeUserApp(c *gin.Context) {
	userID, _ := getUserID(c)

	var req dto.ChangeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, 400, "validation_error", err.Error())
		return
	}

	err := h.services.UpdateUserApplication(c.Request.Context(), userID, req)
	if err != nil {
		writeErr(c, 500, "internal_error", "ошибка изменения")
		return
	}

	c.JSON(200, gin.H{"result": "updated"})
}
