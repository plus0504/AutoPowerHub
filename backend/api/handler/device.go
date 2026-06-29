package handler

import (
	"errors"
	"net/http"
	"strconv"

	devicesvc "autopowerhub/service/device"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	svc *devicesvc.Service
}

func NewDeviceHandler(svc *devicesvc.Service) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	devices, err := h.svc.ListDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func (h *DeviceHandler) Power(c *gin.Context) {
	h.handleCommand(c, "PRESS")
}

func (h *DeviceHandler) Test(c *gin.Context) {
	h.handleCommand(c, "TEST")
}

func (h *DeviceHandler) handleCommand(c *gin.Context, command string) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	username, _ := c.Get("username")

	if err := h.svc.SendCommand(uint(id), username.(string), command); err != nil {
		switch {
		case errors.Is(err, devicesvc.ErrDeviceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, devicesvc.ErrDeviceDisabled):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "command sent"})
}
