package handler

import (
	"errors"
	"net/http"

	cozylifesvc "autopowerhub/service/cozylife"

	"github.com/gin-gonic/gin"
)

type CozyLifeHandler struct {
	svc *cozylifesvc.Service
}

func NewCozyLifeHandler(svc *cozylifesvc.Service) *CozyLifeHandler {
	return &CozyLifeHandler{svc: svc}
}

func (h *CozyLifeHandler) ListSwitches(c *gin.Context) {
	switches, err := h.svc.ListSwitches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if switches == nil {
		switches = []*cozylifesvc.ScanResult{}
	}
	c.JSON(http.StatusOK, gin.H{"switches": switches})
}

func (h *CozyLifeHandler) TurnOn(c *gin.Context) {
	h.handleSwitchCmd(c, h.svc.TurnOn)
}

func (h *CozyLifeHandler) TurnOff(c *gin.Context) {
	h.handleSwitchCmd(c, h.svc.TurnOff)
}

func (h *CozyLifeHandler) handleSwitchCmd(c *gin.Context, fn func(string) error) {
	ip := c.Param("ip")
	if err := fn(ip); err != nil {
		if errors.Is(err, cozylifesvc.ErrSwitchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
