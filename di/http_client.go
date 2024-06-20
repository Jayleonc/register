package di

import (
	"Jayleonc/gateway/pkg/netx"
	"time"
)

func InitHTTPClient() netx.HTTPClientI {
	return netx.NewHTTPClient(10 * time.Second)
}
