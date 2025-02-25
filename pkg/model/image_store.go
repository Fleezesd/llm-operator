package model

import (
	"context"

	llmv1alpha1 "github.com/fleezesd/llm-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ImageStorePVCName         = "ollama-models-store-pvc"
	ImageStoreStatefulsetName = "ollama-models-store"
)

func EnsureImageStorePVCCreated(
	ctx context.Context,
	namespace string,
	storageClassName string,
	pvcSource *corev1.PersistentVolumeClaimVolumeSource,
	pvSpec *llmv1alpha1.ModelPersistentVolumeSpec,
) (*corev1.PersistentVolumeClaim, error) {
	log := log.FromContext(ctx)
	client := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*llmv1alpha1.Model](ctx)

	pvc, err := getImageStorePVC(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	if pvc != nil {
		return pvc, nil
	}

	log.Info("no existing image storage PVC found, creating one...")

	accessMode := corev1.ReadWriteOnce
	if pvSpec != nil && pvSpec.AccessMode != nil {
		accessMode = *pvSpec.AccessMode
	}

	pvc = &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ImageStorePVCName,
			Namespace: namespace,
			Labels:    ImageStoreLabels(),
			// todo: add annonations
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("100Gi"),
				},
			},
			StorageClassName: &storageClassName,
			AccessModes:      []corev1.PersistentVolumeAccessMode{accessMode},
		},
	}
	err = client.Create(ctx, pvc)
	if err != nil {
		return nil, err
	}
	log.Info("created image storage PVC", "pvc", pvc)
	modelRecorder.Event(corev1.EventTypeNormal, "ProvisionedImageStoragePVC", "Provisioned image storage PVC")
	return pvc, nil
}

func getImageStorePVC(ctx context.Context, client client.Client, namespace string) (*corev1.PersistentVolumeClaim, error) {
	var pvc corev1.PersistentVolumeClaim
	err := client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      ImageStorePVCName,
	}, &pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &pvc, err
}

func ImageStoreLabels() map[string]string {
	return map[string]string{
		"app":                "ollama-image-store",
		"ollama.fleezesd.io": "iamge-store",
	}
}
