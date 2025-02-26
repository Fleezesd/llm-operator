package model

import (
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	OllamaBaseImage = "ollama/ollama"
)

func NewOllamaServerContainer(
	readOnly bool,
	resources corev1.ResourceRequirements,
	extraEnvFrom []corev1.EnvFromSource,
	extraEnv []corev1.EnvVar,
) corev1.Container {
	return corev1.Container{
		Name:  "ollama-server",
		Image: OllamaBaseImage,
		Args:  []string{"serve"},
		Ports: []corev1.ContainerPort{
			{
				Name:          "ollama",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: 11434,
			},
		},
		EnvFrom: extraEnvFrom,
		Env: UniqEnvVar(
			append(append([]corev1.EnvVar{}, corev1.EnvVar{
				Name:  "OLLAMA_HOST",
				Value: "0.0.0.0",
			}), extraEnv...),
		),
		Resources: resources,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "image-storage",
				MountPath: "/root/.ollama",
				ReadOnly:  readOnly,
			},
		},
		VolumeDevices: []corev1.VolumeDevice{},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 5,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      1,
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 5,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      5,
		},
	}
}

func UniqEnvVar(env []corev1.EnvVar) []corev1.EnvVar {
	return lo.UniqBy(env, func(item corev1.EnvVar) string {
		return item.Name
	})
}
