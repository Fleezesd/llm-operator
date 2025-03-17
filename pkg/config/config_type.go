package config

import (
	corev1 "k8s.io/api/core/v1"
)

type Config struct {
	// SystemDatasource specifies the built-in datasource to host data files and model files
	SystemDataSource corev1.TypedObjectReference `json:"systemDataSource"`

	// RelationalDatasource specifies the built-in datasource(common:postgres) to host relational data
	RelationDataSource corev1.TypedObjectReference `json:"relationDataSource"`
}
