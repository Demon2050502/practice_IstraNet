package handler

import (
	"net/http"

	"practice_IstraNet/pkg/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminCreateStatus(c *gin.Context) {
	if _, ok := getAdminID(c); !ok {
		return
	}

	var req dto.AdminCreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	resp, err := h.services.CreateStatus(c.Request.Context(), req)
	if err != nil {
		writeAdminStatusErr(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AdminUpdateStatus(c *gin.Context) {
	if _, ok := getAdminID(c); !ok {
		return
	}

	var req dto.AdminUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	resp, err := h.services.UpdateStatus(c.Request.Context(), req)
	if err != nil {
		writeAdminStatusErr(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AdminDeleteStatus(c *gin.Context) {
	if _, ok := getAdminID(c); !ok {
		return
	}

	var req dto.AdminDeleteStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.services.DeleteStatus(c.Request.Context(), req.ID); err != nil {
		writeAdminStatusErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}
