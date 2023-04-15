package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
)

type userIdCtxKey struct{}

func userIdToContext(ctx context.Context, userid int32) context.Context {
	return context.WithValue(ctx, userIdCtxKey{}, userid)
}

func userIdFromContext(ctx context.Context) (int32, bool) {
	uid, ok := ctx.Value(userIdCtxKey{}).(int32)
	if !ok {
		return 0, ok
	}
	return uid, ok
}

