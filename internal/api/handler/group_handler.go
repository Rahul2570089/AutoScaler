package handler

import (
	"AutoScaler/internal/models"
	"AutoScaler/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	storage storage.Storage
}

func NewGroupHandler(storage storage.Storage) *GroupHandler {
	return &GroupHandler{
		storage: storage,
	}
}

type CreateGroupRequest struct {
	Name                   string            `json:"name" binding:"required"`
	LaunchSpec             models.LaunchSpec `json:"launchSpec" binding:"required"`
	MinCapacity            int               `json:"minCapacity" binding:"required,min=0"`
	MaxCapacity            int               `json:"maxCapacity" binding:"required,min=0"`
	DesiredCapacity        int               `json:"desiredCapacity" binding:"required,min=0"`
	Zones                  []string          `json:"zones" binding:"required"`
	HealthCheckType        string            `json:"healthCheckType"`
	HealthCheckGracePeriod int64             `json:"healthCheckGracePeriod"`
	DefaultCooldown        int64             `json:"defaultCooldown"`
	TerminationPolicies    []string          `json:"terminationPolicies"`
}

type UpdateGroupRequest struct {
	LaunchSpec             *models.LaunchSpec `json:"launchSpec"`
	MinCapacity            *int               `json:"minCapacity"`
	MaxCapacity            *int               `json:"maxCapacity"`
	DesiredCapacity        *int               `json:"desiredCapacity"`
	Zones                  []string           `json:"zones"`
	HealthCheckType        *string            `json:"healthCheckType"`
	HealthCheckGracePeriod *int64             `json:"healthCheckGracePeriod"`
	DefaultCooldown        *int64             `json:"defaultCooldown"`
	TerminationPolicies    []string           `json:"terminationPolicies"`
	Version                int64              `json:"version" binding:"required"`
}

type SetDesiredCapacityRequest struct {
	DesiredCapacity int `json:"desiredCapacity" binding:"required,min=0"`
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate min <= desired <= max
	if req.MinCapacity > req.DesiredCapacity || req.DesiredCapacity > req.MaxCapacity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "minCapacity <= desiredCapacity <= maxCapacity required",
		})
		return
	}

	group := &models.Group{
		Name:                   req.Name,
		LaunchSpec:             req.LaunchSpec,
		MinCapacity:            req.MinCapacity,
		MaxCapacity:            req.MaxCapacity,
		DesiredCapacity:        req.DesiredCapacity,
		Zones:                  req.Zones,
		HealthCheckType:        req.HealthCheckType,
		HealthCheckGracePeriod: req.HealthCheckGracePeriod,
		DefaultCooldown:        req.DefaultCooldown,
		TerminationPolicies:    req.TerminationPolicies,
		Version:                1, // Initial version
	}

	if err := h.storage.CreateGroup(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	name := c.Param("name")

	group, err := h.storage.GetGroup(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"groups": group})
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	groups, err := h.storage.ListGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	name := c.Param("name")

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.storage.GetGroup(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if req.LaunchSpec != nil {
		group.LaunchSpec = *req.LaunchSpec
	}
	if req.MinCapacity != nil {
		group.MinCapacity = *req.MinCapacity
	}
	if req.MaxCapacity != nil {
		group.MaxCapacity = *req.MaxCapacity
	}
	if req.DesiredCapacity != nil {
		group.DesiredCapacity = *req.DesiredCapacity
	}
	if req.Zones != nil {
		group.Zones = req.Zones
	}
	if req.HealthCheckType != nil {
		group.HealthCheckType = *req.HealthCheckType
	}
	if req.HealthCheckGracePeriod != nil {
		group.HealthCheckGracePeriod = *req.HealthCheckGracePeriod
	}
	if req.DefaultCooldown != nil {
		group.DefaultCooldown = *req.DefaultCooldown
	}
	if req.TerminationPolicies != nil {
		group.TerminationPolicies = req.TerminationPolicies
	}

	group.Version = req.Version + 1

	if err := h.storage.UpdateGroup(c.Request.Context(), group); err != nil {
		if err.Error() == "stale update" {
			c.JSON(http.StatusConflict, gin.H{"error": "version conflict"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	name := c.Param("name")

	if err := h.storage.DeleteGroup(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *GroupHandler) SetDesiredCapacity(c *gin.Context) {
	name := c.Param("name")

	var req SetDesiredCapacityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.storage.GetGroup(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Validate against min/max
	if req.DesiredCapacity < group.MinCapacity || req.DesiredCapacity > group.MaxCapacity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "desiredCapacity must be between minCapacity and maxCapacity",
		})
		return
	}

	group.DesiredCapacity = req.DesiredCapacity
	group.Version++

	if err := h.storage.UpdateGroup(c.Request.Context(), group); err != nil {
		if err.Error() == "stale update" {
			c.JSON(http.StatusConflict, gin.H{"error": "version conflict"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}
