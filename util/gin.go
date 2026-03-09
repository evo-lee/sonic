package util

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/web/ginadapter"
)

func GetClientIP(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return ginCtx.ClientIP()
}

func GetUserAgent(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return ginCtx.GetHeader("User-Agent")
}

func MustGetQueryString(ctx *gin.Context, key string) (string, error) {
	return MustGetWebQueryString(ginadapter.NewContext(ctx), key)
}

func MustGetQueryInt32(ctx *gin.Context, key string) (int32, error) {
	return MustGetWebQueryInt32(ginadapter.NewContext(ctx), key)
}

func MustGetQueryInt64(ctx *gin.Context, key string) (int64, error) {
	return MustGetWebQueryInt64(ginadapter.NewContext(ctx), key)
}

func MustGetQueryInt(ctx *gin.Context, key string) (int, error) {
	return MustGetWebQueryInt(ginadapter.NewContext(ctx), key)
}

func MustGetQueryBool(ctx *gin.Context, key string) (bool, error) {
	return MustGetWebQueryBool(ginadapter.NewContext(ctx), key)
}

func GetQueryBool(ctx *gin.Context, key string, defaultValue bool) (bool, error) {
	return GetWebQueryBool(ginadapter.NewContext(ctx), key, defaultValue)
}

func ParamString(ctx *gin.Context, key string) (string, error) {
	return ParamWebString(ginadapter.NewContext(ctx), key)
}

func ParamInt32(ctx *gin.Context, key string) (int32, error) {
	return ParamWebInt32(ginadapter.NewContext(ctx), key)
}

func ParamInt64(ctx *gin.Context, key string) (int64, error) {
	return ParamWebInt64(ginadapter.NewContext(ctx), key)
}

func ParamBool(ctx *gin.Context, key string) (bool, error) {
	return ParamWebBool(ginadapter.NewContext(ctx), key)
}
