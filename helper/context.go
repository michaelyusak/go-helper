package helper

import (
	"context"

	"github.com/michaelyusak/go-helper/appconstant"
)

func InjectValues(ctx context.Context, values map[appconstant.ContextKey]any) context.Context {
	for k, v := range values {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx
}
