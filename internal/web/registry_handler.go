package web

import (
	"Jayleonc/register/internal/domain"
	"Jayleonc/register/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type RegistryHandler struct {
	svc      domain.Registry
	resolver domain.Resolver
}

func NewRegistryHandler(svc domain.Registry, resolver domain.Resolver) *RegistryHandler {
	return &RegistryHandler{
		svc:      svc,
		resolver: resolver,
	}
}

func (h *RegistryHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/registry")
	rg.POST("/register", h.Register)
	rg.POST("/unregister", h.Unregister)
	rg.GET("/services/:appId", h.Discover)
	rg.GET("/health", h.HealthCheck)
}

func (h *RegistryHandler) Register(ctx *gin.Context) {
	var instance domain.ServiceInstance
	requestID := uuid.New().String()

	if err := ctx.BindJSON(&instance); err != nil {
		ginx.Error(ctx, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := h.svc.Register(instance); err != nil {
		ginx.Error(ctx, http.StatusInternalServerError, "注册失败", err)
		return
	}

	ginx.OK(ctx, ginx.Response{Msg: "注册成功", Data: requestID})
}

func (h *RegistryHandler) Unregister(ctx *gin.Context) {
	var request struct {
		AppID      string `json:"appId"`
		InstanceID string `json:"instanceId"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		ginx.Error(ctx, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := h.svc.Unregister(request.AppID, request.InstanceID); err != nil {
		ginx.Error(ctx, http.StatusInternalServerError, "注销失败", err)
		return
	}

	ginx.OK(ctx, ginx.Response{Msg: "注销成功"})
}

func (h *RegistryHandler) Discover(ctx *gin.Context) {
	appID := ctx.Param("appId")
	instances, err := h.resolver.Resolve(appID)
	if err != nil {
		ginx.Error(ctx, http.StatusNotFound, "服务未找到", err)
		return
	}
	ginx.OK(ctx, ginx.Response{Msg: "查询成功", Data: instances})
}

func (h *RegistryHandler) HealthCheck(ctx *gin.Context) {
	if err := h.svc.HealthCheck(); err != nil {
		ginx.Error(ctx, http.StatusInternalServerError, "健康检查失败", err)
		return
	}
	ginx.OK(ctx, ginx.Response{Msg: "健康检查通过"})
}
