package web

import (
	"context"
	"github.com/Jayleonc/register/config_center"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ConfigHandler struct {
	configCenter *config_center.Client
}

func NewConfigHandler(configCenter *config_center.Client) *ConfigHandler {
	return &ConfigHandler{
		configCenter: configCenter,
	}
}

func (h *ConfigHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/put", h.putConfigHandler)
	r.GET("/get", h.getConfigHandler)
	r.DELETE("/delete", h.deleteConfigHandler)
	r.GET("/watch", h.watchConfigHandler)
	r.GET("/list", h.listConfigHandler) // 添加新的路由
}

func (h *ConfigHandler) putConfigHandler(c *gin.Context) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.configCenter.PutConfig(context.Background(), req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ConfigHandler) getConfigHandler(c *gin.Context) {
	key := c.Query("key")

	value, err := h.configCenter.GetConfig(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func (h *ConfigHandler) deleteConfigHandler(c *gin.Context) {
	key := c.Query("key")

	if err := h.configCenter.DeleteConfig(context.Background(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ConfigHandler) watchConfigHandler(c *gin.Context) {
	key := c.Query("key")

	ch, err := h.configCenter.WatchConfig(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		if value, ok := <-ch; ok {
			c.SSEvent("message", gin.H{"key": key, "value": value})
			return true
		}
		return false
	})
}

func (h *ConfigHandler) listConfigHandler(c *gin.Context) {
	prefix := c.Query("prefix")

	configs, err := h.configCenter.ListConfig(context.Background(), prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}
