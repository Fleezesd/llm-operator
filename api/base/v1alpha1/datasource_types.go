/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DataSourceSpec defines the desired state of DataSource
type DataSourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	CommonSpec `json:",inline"`

	// Endpoint defines connection info
	Endpoint Endpoint `json:"endpoint"`

	// OSS defines info for object storage service
	OSS *OSS `json:"oss,omitempty"`

	// RDMA configure RDMA pulls the model file directly from the remote service to the host node.
	RDMA *RDMA `json:"rdma,omitempty"`

	// PostgreSQL defines info for PostgreSQL
	PostgreSQL *PostgreSQL `json:"postgresql,omitempty"`

	// Web defines info for web resources
	Web *Web `json:"web,omitempty"`
}

// OSS defines info for object storage service as datasource
type OSS struct {
	Bucket string `json:"bucket,omitempty"`

	Object string `json:"object,omitempty"`

	VersionID string `json:"versionID,omitempty"`
}

type RDMA struct {
	// We consider the model storage path on the sender's side and the save path on the receiver's side to be the same,
	// so a single Path is uniformly configured here.
	// example: /opt/kubeagi/, /opt/, /
	// +kubebuilder:validation:Pattern=(^\/$)|(^\/[a-zA-Z0-9\_.@-]+(\/[a-zA-Z0-9\_.@-]+)*\/$)
	Path string `json:"path"`

	NodePaths map[string]string `json:"nodePaths,omitempty"`
}

// PostgreSQL defines info for PostgreSQL
//
// ref: https://github.com/jackc/pgx/blame/v5.5.1/pgconn/config.go#L409
// they are common standard PostgreSQL environment variables
// For convenience, we use the same name.
//
// The PGUSER/PGPASSWORD/PGPASSFILE/PGSSLPASSWORD parameters have been intentionally excluded
// because they contain sensitive information and are stored in the secret pointed to by `endpoint.authSecret`.
type PostgreSQL struct {
	Host               string `json:"PGHOST,omitempty"`
	Port               string `json:"PGPORT,omitempty"`
	Database           string `json:"PGDATABASE,omitempty"`
	AppName            string `json:"PGAPPNAME,omitempty"`
	ConnectTimeout     string `json:"PGCONNECT_TIMEOUT,omitempty"`
	SSLMode            string `json:"PGSSLMODE,omitempty"`
	SSLKey             string `json:"PGSSLKEY,omitempty"`
	SSLCert            string `json:"PGSSLCERT,omitempty"`
	SSLSni             string `json:"PGSSLSNI,omitempty"`
	SSLRootCert        string `json:"PGSSLROOTCERT,omitempty"`
	TargetSessionAttrs string `json:"PGTARGETSESSIONATTRS,omitempty"`
	Service            string `json:"PGSERVICE,omitempty"`
	ServiceFile        string `json:"PGSERVICEFILE,omitempty"`
}

const (
	PGUSER        = "PGUSER"
	PGPASSWORD    = "PGPASSWORD"
	PGPASSFILE    = "PGPASSFILE"
	PGSSLPASSWORD = "PGSSLPASSWORD"
)

// Web defines info for web resources
type Web struct {
	// RecommendIntervalTime is the recommended interval time for this crawler
	RecommendIntervalTime int `json:"recommendIntervalTime,omitempty"`
}

// DataSourceStatus defines the observed state of DataSource
type DataSourceStatus struct {
	// ConditionedStatus is the current status
	ConditionedStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced
//+kubebuilder:printcolumn:name="display-name",type=string,JSONPath=`.spec.displayName`
//+kubebuilder:printcolumn:name="type",type=string,JSONPath=`.metadata.labels.fleezesd\.k8s\.com\.cn/datasource-type`

// DataSource is the Schema for the datasources API
type DataSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataSourceSpec   `json:"spec,omitempty"`
	Status DataSourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataSourceList contains a list of DataSource
type DataSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataSource{}, &DataSourceList{})
}
