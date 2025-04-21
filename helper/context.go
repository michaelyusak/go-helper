package helper

import "context"

func InjectValues(ctx context.Context, values map[any]any) context.Context {
	for k, v := range values {
		ctx = context.WithValue(ctx, k, v)
	}
	
	return ctx
}
