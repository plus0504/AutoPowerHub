package handler

import (
	"net/http"
	"strconv"

	"autopowerhub/repository"
	"autopowerhub/service/ble"

	"github.com/gin-gonic/gin"
)

type DebugHandler struct {
	deviceRepo *repository.DeviceRepository
	bleMgr     *ble.Manager
}

func NewDebugHandler(deviceRepo *repository.DeviceRepository, bleMgr *ble.Manager) *DebugHandler {
	return &DebugHandler{deviceRepo: deviceRepo, bleMgr: bleMgr}
}

// BLEScan connects to the device and dumps all BLE service/characteristic UUIDs.
// Useful for diagnosing UUID mismatches between firmware and config.yaml.
func (h *DebugHandler) BLEScan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	device, err := h.deviceRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	results, err := h.bleMgr.DumpServices(device.MAC)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mac":      device.MAC,
		"services": results,
	})
}
