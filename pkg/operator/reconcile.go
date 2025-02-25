package operator

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Reconcilers[T runtime.Object] []SubReconciler[T]

func NewReconcilers[T runtime.Object](reconcilers ...SubReconciler[T]) Reconcilers[T] {
	return reconcilers
}

// this func can not use ptr
func (s Reconcilers[T]) Reconcile(ctx context.Context, req ctrl.Request, obj T) error {
	log := log.FromContext(ctx)

	for _, subReconciler := range s {
		log.Info("Reconciling", "kind", subReconciler.kind)

		err := subReconciler.reconcile(ctx, req.Namespace, req.Name, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

type ReconcileHandler[T runtime.Object] func(ctx context.Context, namespace string, name string, obj T) error
type SubReconciler[T runtime.Object] struct {
	apiverison string
	group      string
	kind       string
	reconcile  func(ctx context.Context, namespace string, name string, obj T) error
}

func NewPVCReconciler[T runtime.Object](fn ReconcileHandler[T]) SubReconciler[T] {
	return SubReconciler[T]{
		apiverison: "v1",
		group:      "core",
		kind:       "PersistentVolumeClaim",
		reconcile:  fn,
	}
}
