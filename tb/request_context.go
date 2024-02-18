package tb

import (
	"context"
	"net/http"
)

type contextKey int

const (
	Attempts contextKey = iota
	Retry
)

func GetRetryFromContext(request *http.Request) int {
	if retry, ok := request.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

func SetRetryForContext(request *http.Request, retries int) context.Context {
	ctx := context.WithValue(request.Context(), Retry, retries+1)
	return ctx
}

func GetAttemptsFromContext(request *http.Request) int {
	if attempts, ok := request.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func SetAttemptsForContext(request *http.Request, attempts int) context.Context {
	ctx := context.WithValue(request.Context(), Attempts, attempts)
	return ctx
}
