package registry

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Param struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type Return struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type Api struct {
	Method  string   `json:"method"`
	Path    string   `json:"path"`
	Params  []Param  `json:"params"`
	Returns []Return `json:"returns"`
}

type App struct {
	mu         sync.Mutex
	interfaces []Api
	*gin.Engine
}

func NewApiDescriptor(engine *gin.Engine) *App {
	return &App{
		interfaces: make([]Api, 0),
		Engine:     engine,
	}
}

func (b *App) AddInterface(method, path string, params []Param, returns []Return) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.interfaces = append(b.interfaces, Api{
		Method:  method,
		Path:    path,
		Params:  params,
		Returns: returns,
	})
}

func (b *App) GetInterfaces() []Api {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.interfaces
}

func (b *App) GET(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("GET", path, params, returns)
	b.Engine.GET(path, handler)
}

func (b *App) POST(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("POST", path, params, returns)
	b.Engine.POST(path, handler)
}

func (b *App) PUT(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("PUT", path, params, returns)
	b.Engine.PUT(path, handler)
}

func (b *App) DELETE(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("DELETE", path, params, returns)
	b.Engine.DELETE(path, handler)
}

// Group 新增 Group 方法
func (b *App) Group(relativePath string, handlers ...gin.HandlerFunc) *ApiGroup {
	group := b.Engine.Group(relativePath, handlers...)
	return &ApiGroup{
		App:         b,
		RouterGroup: group,
	}
}

type ApiGroup struct {
	App *App
	*gin.RouterGroup
}

func (g *ApiGroup) GET(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	fullPath := g.calculateFullPath(path)
	g.App.AddInterface("GET", fullPath, params, returns)
	g.RouterGroup.GET(path, handler)
}

func (g *ApiGroup) POST(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	fullPath := g.calculateFullPath(path)
	g.App.AddInterface("POST", fullPath, params, returns)
	g.RouterGroup.POST(path, handler)
}

func (g *ApiGroup) PUT(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	fullPath := g.calculateFullPath(path)
	g.App.AddInterface("PUT", fullPath, params, returns)
	g.RouterGroup.PUT(path, handler)
}

func (g *ApiGroup) DELETE(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	fullPath := g.calculateFullPath(path)
	g.App.AddInterface("DELETE", fullPath, params, returns)
	g.RouterGroup.DELETE(path, handler)
}

func (g *ApiGroup) calculateFullPath(relativePath string) string {
	return g.BasePath() + relativePath
}
