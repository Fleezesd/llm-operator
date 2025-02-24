package model

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type WrappedRecorder[T runtime.Object] struct {
	recorder record.EventRecorder
	t        T
}

func NewWrappedRecorder[T runtime.Object](recorder record.EventRecorder, object T) *WrappedRecorder[T] {
	return &WrappedRecorder[T]{
		recorder: recorder,
		t:        object,
	}
}

type baseWrapperRecorderContextKey string

const (
	defaultBaseWrapperRecorderContextKey baseWrapperRecorderContextKey = "default"
)

func WithWrappedRecorder[T runtime.Object](
	ctx context.Context,
	recorder *WrappedRecorder[T],
	key ...baseWrapperRecorderContextKey,
) context.Context {
	if len(key) == 0 {
		return context.WithValue(ctx, defaultBaseWrapperRecorderContextKey, recorder)
	}
	return context.WithValue(ctx, key[0], recorder)
}

func WrappedRecorderFromContext[T runtime.Object](
	ctx context.Context,
	key ...baseWrapperRecorderContextKey,
) *WrappedRecorder[T] {
	if len(key) == 0 {
		// default wrapped Recorder
		r, _ := ctx.Value(defaultBaseWrapperRecorderContextKey).(*WrappedRecorder[T])
		return r
	}
	r, _ := ctx.Value(key[0]).(*WrappedRecorder[T])
	return r
}
