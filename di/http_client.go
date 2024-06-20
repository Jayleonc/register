package di

import (
	"Jayleonc/register/pkg/netx"
	"time"
)

func InitHTTPClient() netx.HTTPClientI {
	return netx.NewHTTPClient(10 * time.Second)
}
