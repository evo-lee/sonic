package util

import "context"

func GetClientIP(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if clientIP, ok := ctx.Value("client_ip").(string); ok {
		return clientIP
	}
	return ""
}

func GetUserAgent(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userAgent, ok := ctx.Value("user_agent").(string); ok {
		return userAgent
	}
	return ""
}
