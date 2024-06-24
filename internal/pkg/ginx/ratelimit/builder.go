package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/Jayleonc/register/internal/pkg/limiter"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
