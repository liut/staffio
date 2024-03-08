package backends

import (
	"context"
)

type contextKey int

const (
	siteK contextKey = iota
)

func ContextWithSite(ctx context.Context, site string) context.Context {
	if len(site) == 0 {
		return ctx
	}

	return context.WithValue(ctx, siteK, site)
}

func SiteFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(siteK).(string); ok {
		return val
	}
	return ""
}
