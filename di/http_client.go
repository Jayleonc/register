package di

import (
	"git.daochat.cn/service/registry/pkg/netx"
	"time"
)

func InitHTTPClient() netx.HTTPClientI {
	return netx.NewHTTPClient(10 * time.Second)
}
