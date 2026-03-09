package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/web/ginadapter"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type AuthMiddleware struct {
	OptionService       service.OptionService
	OneTimeTokenService service.OneTimeTokenService
	UserService         service.UserService
	Cache               cache.Cache
}

func NewAuthMiddleware(optionService service.OptionService, oneTimeTokenService service.OneTimeTokenService, cache cache.Cache, userService service.UserService) *AuthMiddleware {
	authMiddleware := &AuthMiddleware{
		OptionService:       optionService,
		OneTimeTokenService: oneTimeTokenService,
		Cache:               cache,
		UserService:         userService,
	}
	return authMiddleware
}

func (a *AuthMiddleware) GetWrapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		webCtx := ginadapter.NewContext(ctx)
		isInstalled, err := a.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
		if err != nil {
			abortWithStatusJSON(webCtx, http.StatusInternalServerError, "")
			return
		}
		if !isInstalled.(bool) {
			abortWithStatusJSON(webCtx, http.StatusBadRequest, T(webCtx, "auth.blog_not_initialized", "Blog is not initialized"))
			return
		}

		oneTimeToken, ok := ctx.GetQuery(consts.OneTimeTokenQueryName)
		if ok {
			allowedURL, ok := a.OneTimeTokenService.Get(oneTimeToken)
			if !ok {
				abortWithStatusJSON(webCtx, http.StatusBadRequest, T(webCtx, "auth.one_time_token_not_exist_or_expired", "OneTimeToken is not exist or expired"))
				return
			}
			currentURL := ctx.Request.URL.Path
			if currentURL != allowedURL {
				abortWithStatusJSON(webCtx, http.StatusBadRequest, T(webCtx, "auth.one_time_token_uri_mismatch", "The one-time token does not correspond the request uri"))
				return
			}
			return
		}

		token := ctx.GetHeader(consts.AdminTokenHeaderName)
		if token == "" {
			abortWithStatusJSON(webCtx, http.StatusUnauthorized, T(webCtx, "auth.not_logged_in", "Not logged in, please login first"))
			return
		}
		userID, ok := a.Cache.Get(cache.BuildTokenAccessKey(token))

		if !ok || userID == nil {
			abortWithStatusJSON(webCtx, http.StatusUnauthorized, T(webCtx, "auth.token_expired_or_not_exist", "Token has expired or does not exist"))
			return
		}

		user, err := a.UserService.GetByID(ctx, userID.(int32))
		if xerr.GetType(err) == xerr.NoRecord {
			_ = ctx.Error(err)
			abortWithStatusJSON(webCtx, http.StatusUnauthorized, T(webCtx, "auth.user_not_found", "User not found"))
			return
		}
		if err != nil {
			_ = ctx.Error(err)
			abortWithStatusJSON(webCtx, http.StatusInternalServerError, "")
			return
		}
		ctx.Set(consts.AuthorizedUser, user)
	}
}
