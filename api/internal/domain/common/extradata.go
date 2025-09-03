package common

import "context"

type ExtraData map[string]any
type ContextKey string

const extraDataCtxKey ContextKey = "extraData"

// Sets `ExtraData` in the given context
func SetExtraDataInCtx(ctx context.Context, extraData ExtraData) context.Context {
	return context.WithValue(ctx, extraDataCtxKey, extraData)
}

// Extracts `ExtraData` from the given context
func GetExtraDataFromCtx(ctx context.Context) ExtraData {
	ctxVal := ctx.Value(extraDataCtxKey)
	switch ex := ctxVal.(type) {
	case ExtraData:
		return ex
	default:
		return make(ExtraData)
	}
}
