package config

import (
	"context"
	"errors"

	basev1alpha1 "github.com/fleezesd/llm-operator/api/base/v1alpha1"
)

var (
	ErrSystemCliNotFound = errors.New("system cli not found")
)

func GetSystemDatasource(ctx context.Context) *basev1alpha1.DataSource {
	return nil
}

func getConfig(ctx context.Context) (config *Config, err error) {
	return nil, nil
}
