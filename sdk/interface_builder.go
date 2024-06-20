package sdk

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Return struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ServiceInterface struct {
	Method  string   `json:"method"`
	Path    string   `json:"path"`
	Params  []Param  `json:"params"`
	Returns []Return `json:"returns"`
}

type InterfaceBuilder struct {
	mu         sync.Mutex
	interfaces []ServiceInterface
}

func NewInterfaceBuilder() *InterfaceBuilder {
	return &InterfaceBuilder{
		interfaces: make([]ServiceInterface, 0),
	}
}

func (b *InterfaceBuilder) AddInterface(method, path string, params []Param, returns []Return) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.interfaces = append(b.interfaces, ServiceInterface{
		Method:  method,
		Path:    path,
		Params:  params,
		Returns: returns,
	})
}

func (b *InterfaceBuilder) GetInterfaces() []ServiceInterface {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.interfaces
}

type HandlerBuilder struct {
	interfaceBuilder *InterfaceBuilder
	method           string
	path             string
	params           []Param
	returns          []Return
}

func NewHandlerBuilder(interfaceBuilder *InterfaceBuilder, method, path string) *HandlerBuilder {
	return &HandlerBuilder{
		interfaceBuilder: interfaceBuilder,
		method:           method,
		path:             path,
		params:           []Param{},
		returns:          []Return{},
	}
}

func (hb *HandlerBuilder) AddParam(name, paramType string) *HandlerBuilder {
	hb.params = append(hb.params, Param{Name: name, Type: paramType})
	return hb
}

func (hb *HandlerBuilder) AddReturn(name, returnType string) *HandlerBuilder {
	hb.returns = append(hb.returns, Return{Name: name, Type: returnType})
	return hb
}

func (hb *HandlerBuilder) Build(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	hb.interfaceBuilder.AddInterface(hb.method, hb.path, hb.params, hb.returns)
	return handlerFunc
}
