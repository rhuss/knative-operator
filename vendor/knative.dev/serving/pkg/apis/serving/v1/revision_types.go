/*
Copyright 2019 The Knative Authors

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Revision is an immutable snapshot of code and configuration.  A revision
// references a container image. Revisions are created by updates to a
// Configuration.
//
// See also: https://github.com/knative/serving/blob/main/docs/spec/overview.md#revision
type Revision struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec RevisionSpec `json:"spec,omitempty"`

	// +optional
	Status RevisionStatus `json:"status,omitempty"`
}

// Verify that Revision adheres to the appropriate interfaces.
var (
	// Check that Revision can be validated, can be defaulted, and has immutable fields.
	_ apis.Validatable = (*Revision)(nil)
	_ apis.Defaultable = (*Revision)(nil)

	// Check that Revision can be converted to higher versions.
	_ apis.Convertible = (*Revision)(nil)

	// Check that we can create OwnerReferences to a Revision.
	_ kmeta.OwnerRefable = (*Revision)(nil)

	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Revision)(nil)
)

// RevisionTemplateSpec describes the data a revision should have when created from a template.
// Based on: https://github.com/kubernetes/api/blob/e771f807/core/v1/types.go#L3179-L3190
type RevisionTemplateSpec struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec RevisionSpec `json:"spec,omitempty"`
}

// RevisionSpec holds the desired state of the Revision (from the client).
type RevisionSpec struct {
	corev1.PodSpec `json:",inline"`

	// ContainerConcurrency specifies the maximum allowed in-flight (concurrent)
	// requests per container of the Revision.  Defaults to `0` which means
	// concurrency to the application is not limited, and the system decides the
	// target concurrency for the autoscaler.
	// +optional
	ContainerConcurrency *int64 `json:"containerConcurrency,omitempty"`

	// TimeoutSeconds is the maximum duration in seconds that the request routing
	// layer will wait for a request delivered to a container to begin replying
	// (send network traffic). If unspecified, a system default will be provided.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// MaxDurationSeconds is the maximum duration in seconds a request will be allowed
	// to stay open.
	// +optional
	MaxDurationSeconds *int64 `json:"maxDurationSeconds,omitempty"`
}

const (
	// RevisionConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	RevisionConditionReady = apis.ConditionReady

	// RevisionConditionResourcesAvailable is set when underlying
	// Kubernetes resources have been provisioned.
	RevisionConditionResourcesAvailable apis.ConditionType = "ResourcesAvailable"

	// RevisionConditionContainerHealthy is set when the revision readiness check completes.
	RevisionConditionContainerHealthy apis.ConditionType = "ContainerHealthy"

	// RevisionConditionActive is set when the revision is receiving traffic.
	RevisionConditionActive apis.ConditionType = "Active"
)

// IsRevisionCondition returns true if the ConditionType is a revision condition type
func IsRevisionCondition(t apis.ConditionType) bool {
	switch t {
	case
		RevisionConditionReady,
		RevisionConditionResourcesAvailable,
		RevisionConditionContainerHealthy,
		RevisionConditionActive:
		return true
	}
	return false
}

// RevisionStatus communicates the observed state of the Revision (from the controller).
type RevisionStatus struct {
	duckv1.Status `json:",inline"`

	// LogURL specifies the generated logging url for this particular revision
	// based on the revision url template specified in the controller's config.
	// +optional
	LogURL string `json:"logUrl,omitempty"`

	// ContainerStatuses is a slice of images present in .Spec.Container[*].Image
	// to their respective digests and their container name.
	// The digests are resolved during the creation of Revision.
	// ContainerStatuses holds the container name and image digests
	// for both serving and non serving containers.
	// ref: http://bit.ly/image-digests
	// +optional
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`

	// InitContainerStatuses is a slice of images present in .Spec.InitContainer[*].Image
	// to their respective digests and their container name.
	// The digests are resolved during the creation of Revision.
	// ContainerStatuses holds the container name and image digests
	// for both serving and non serving containers.
	// ref: http://bit.ly/image-digests
	// +optional
	InitContainerStatuses []ContainerStatus `json:"initContainerStatuses,omitempty"`

	// ActualReplicas reflects the amount of ready pods running this revision.
	// +optional
	ActualReplicas *int32 `json:"actualReplicas,omitempty"`
	// DesiredReplicas reflects the desired amount of pods running this revision.
	// +optional
	DesiredReplicas *int32 `json:"desiredReplicas,omitempty"`
}

// ContainerStatus holds the information of container name and image digest value
type ContainerStatus struct {
	Name        string `json:"name,omitempty"`
	ImageDigest string `json:"imageDigest,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RevisionList is a list of Revision resources
type RevisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Revision `json:"items"`
}

// GetStatus retrieves the status of the Revision. Implements the KRShaped interface.
func (t *Revision) GetStatus() *duckv1.Status {
	return &t.Status.Status
}
