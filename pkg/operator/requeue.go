package operator

import (
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type RequeueError struct {
	after time.Duration
	err   error
}

func (r *RequeueError) Error() string {
	if r.err != nil {
		return fmt.Sprintf("encountered an error: %v, requeue after %d", r.err, r.after.Milliseconds())
	}
	return fmt.Sprintf("requeue after %d", r.after.Milliseconds())
}

func (r *RequeueError) Result() *reconcile.Result {
	return &reconcile.Result{Requeue: true, RequeueAfter: r.after}
}
