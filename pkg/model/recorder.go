package model

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// WrappedRecorder wraps a recorder and runtime.object
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

// The resulting event will be created in the same namespace as the reference object.
func (r *WrappedRecorder[T]) Event(eventType, reason, message string) {
	r.recorder.Event(r.t, eventType, reason, message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func (r *WrappedRecorder[T]) Eventf(eventType, reason, messageFmt string, args ...any) {
	r.recorder.Eventf(r.t, eventType, reason, messageFmt, args...)
}

// AnnotatedEventf is just like eventf, but with annotations attached
func (r *WrappedRecorder[T]) AnnotatedEventf(annotations map[string]string, eventtype, reason, messageFmt string, args ...any) {
	r.recorder.AnnotatedEventf(r.t, annotations, eventtype, reason, messageFmt, args...)
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
