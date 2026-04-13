package application

import "context"

type authorizationHeaderContextKey struct{}

func WithAuthorizationHeader(ctx context.Context, header string) context.Context {
	return context.WithValue(ctx, authorizationHeaderContextKey{}, header)
}

func AuthorizationHeaderFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(authorizationHeaderContextKey{}).(string)
	if !ok || value == "" {
		return "", false
	}

	return value, true
}
