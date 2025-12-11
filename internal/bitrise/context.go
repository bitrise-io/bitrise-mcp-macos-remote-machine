package bitrise

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	keyPAT ctxKey = iota
	keyEnabledGroups
)

func patFromCtx(ctx context.Context) (string, error) {
	v := ctx.Value(keyPAT)
	u, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type %T", v)
	}
	return u, nil
}

func ContextWithPAT(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, keyPAT, s)
}
