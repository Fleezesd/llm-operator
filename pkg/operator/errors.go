package operator

import (
	"context"

	"github.com/samber/lo"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func HandleError(ctx context.Context, result *ctrl.Result, err error) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	if result != nil && err != nil {
		log.Error(err, "Requeue", "after", result.RequeueAfter)
		return *result, err
	}
	if result != nil {
		return *result, nil
	}
	return ctrl.Result{}, err
}

func ResultFromError(err error) (*ctrl.Result, error) {
	if lo.IsNil(err) {
		return &ctrl.Result{}, nil
	}
	if requeueErr, ok := err.(*RequeueError); ok {
		return requeueErr.Result(), requeueErr.err
	}
	return nil, err
}
