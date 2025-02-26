package model

import (
	corev1 "k8s.io/api/core/v1"
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
	return corev1.Container{}
}
