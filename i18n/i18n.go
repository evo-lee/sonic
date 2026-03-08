package i18n

import (
	"embed"
	"encoding/json"
	"strings"
)

const (
	LocaleZHCN    = "zh-CN"
	LocaleENUS    = "en-US"
	DefaultLocale = LocaleZHCN
)

var (
	//go:embed locales/*.json
	localeFS embed.FS
	bundles  = map[string]map[string]string{}
)

func init() {
	loadLocale(LocaleZHCN, "locales/zh-CN.json")
	loadLocale(LocaleENUS, "locales/en-US.json")
}

func loadLocale(locale, file string) {
	content, err := localeFS.ReadFile(file)
	if err != nil {
		panic(err)
	}
	temp := map[string]string{}
	if err = json.Unmarshal(content, &temp); err != nil {
		panic(err)
	}
	bundles[locale] = temp
}

func NormalizeLocale(locale string) string {
	locale = strings.TrimSpace(strings.ToLower(strings.ReplaceAll(locale, "_", "-")))
	if locale == "" {
		return ""
	}
	switch {
	case strings.HasPrefix(locale, "zh"):
		return LocaleZHCN
	case strings.HasPrefix(locale, "en"):
		return LocaleENUS
	default:
		return ""
	}
}

func LocaleFromAcceptLanguage(acceptLanguage string) string {
	if acceptLanguage == "" {
		return ""
	}
	parts := strings.Split(acceptLanguage, ",")
	for _, part := range parts {
		token := strings.TrimSpace(strings.Split(part, ";")[0])
		if normalized := NormalizeLocale(token); normalized != "" {
			return normalized
		}
	}
	return ""
}

func ResolveLocale(preferredLocale, acceptLanguage, fallbackLocale string) string {
	if preferred := NormalizeLocale(preferredLocale); preferred != "" {
		return preferred
	}
	if fromBrowser := LocaleFromAcceptLanguage(acceptLanguage); fromBrowser != "" {
		return fromBrowser
	}
	if fallback := NormalizeLocale(fallbackLocale); fallback != "" {
		return fallback
	}
	return DefaultLocale
}

func T(locale, key, fallback string) string {
	locale = ResolveLocale(locale, "", DefaultLocale)
	if value, ok := bundles[locale][key]; ok {
		return value
	}
	if value, ok := bundles[DefaultLocale][key]; ok {
		return value
	}
	if fallback != "" {
		return fallback
	}
	return key
}
