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

type mealIdCtxKey struct{}

func mealIdToContext(ctx context.Context, mealid int32) context.Context {
	return context.WithValue(ctx, mealIdCtxKey{}, mealid)
}

func mealIdFromContext(ctx context.Context) (int32, bool) {
	mealid, ok := ctx.Value(mealIdCtxKey{}).(int32)
	if !ok {
		return 0, ok
	}
	return mealid, ok
}

func mealIdParamToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp := chi.URLParamFromCtx(r.Context(), "mealid")
		if temp == "" {
			exception.WriteErrorAsJSON(w, http.StatusBadRequest, errors.New("no mealid in URL params"))
			return
		}
		t, err := strconv.Atoi(temp)
		if err != nil {
			exception.WriteErrorAsJSON(w, http.StatusBadGateway, errors.New("mealid must be of type int"))
		}
		r = r.WithContext(mealIdToContext(r.Context(), int32(t)))
		next.ServeHTTP(w, r)
	})
}
