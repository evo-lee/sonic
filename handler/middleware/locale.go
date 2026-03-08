package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/i18n"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

const (
	LocaleKey    = "locale"
	LocaleHeader = "Content-Language"
)

type LocaleMiddleware struct {
	optionService service.OptionService
}

func NewLocaleMiddleware(optionService service.OptionService) *LocaleMiddleware {
	return &LocaleMiddleware{optionService: optionService}
}

func (m *LocaleMiddleware) Locale() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		preferred := ""
		options, err := m.optionService.ListAllOption(ctx)
		if err == nil {
			for _, option := range options {
				if option.Key == property.BlogLocale.KeyValue {
					if value, ok := option.Value.(string); ok {
						preferred = value
					}
					break
				}
			}
		}
		locale := i18n.ResolveLocale(preferred, ctx.GetHeader("Accept-Language"), i18n.DefaultLocale)
		ctx.Set(LocaleKey, locale)
		ctx.Writer.Header().Set(LocaleHeader, locale)
		ctx.Next()
	}
}

func GetLocale(ctx *gin.Context) string {
	value, ok := ctx.Get(LocaleKey)
	if !ok {
		return i18n.DefaultLocale
	}
	if locale, ok := value.(string); ok {
		return i18n.ResolveLocale(locale, "", i18n.DefaultLocale)
	}
	return i18n.DefaultLocale
}
