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

type ServiceInterface struct {
	Method  string   `json:"method"`
	Path    string   `json:"path"`
	Params  []Param  `json:"params"`
	Returns []Return `json:"returns"`
}

type ApiDescriptor struct {
	mu         sync.Mutex
	interfaces []ServiceInterface
	engine     *gin.Engine
}

func NewApiDescriptor(engine *gin.Engine) *ApiDescriptor {
	return &ApiDescriptor{
		interfaces: make([]ServiceInterface, 0),
		engine:     engine,
	}
}

func (b *ApiDescriptor) AddInterface(method, path string, params []Param, returns []Return) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.interfaces = append(b.interfaces, ServiceInterface{
		Method:  method,
		Path:    path,
		Params:  params,
		Returns: returns,
	})
}

func (b *ApiDescriptor) GetInterfaces() []ServiceInterface {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.interfaces
}

func (b *ApiDescriptor) GET(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("GET", path, params, returns)
	b.engine.GET(path, handler)
}

func (b *ApiDescriptor) POST(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("POST", path, params, returns)
	b.engine.POST(path, handler)
}

func (b *ApiDescriptor) PUT(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("PUT", path, params, returns)
	b.engine.PUT(path, handler)
}

func (b *ApiDescriptor) DELETE(path string, handler gin.HandlerFunc, params []Param, returns []Return) {
	b.AddInterface("DELETE", path, params, returns)
	b.engine.DELETE(path, handler)
}
