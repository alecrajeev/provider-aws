/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PrivateDNSNamespaceParameters defines the desired state of PrivateDNSNamespace
type PrivateDNSNamespaceParameters struct {
	// Region is which region the PrivateDNSNamespace will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// A description for the namespace.
	Description *string `json:"description,omitempty"`
	// The name that you want to assign to this namespace. When you create a private
	// DNS namespace, AWS Cloud Map automatically creates an Amazon Route 53 private
	// hosted zone that has the same name as the namespace.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// The tags to add to the namespace. Each tag consists of a key and an optional
	// value, both of which you define. Tag keys can have a maximum character length
	// of 128 characters, and tag values can have a maximum length of 256 characters.
	Tags                                []*Tag `json:"tags,omitempty"`
	CustomPrivateDNSNamespaceParameters `json:",inline"`
}

// PrivateDNSNamespaceSpec defines the desired state of PrivateDNSNamespace
type PrivateDNSNamespaceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PrivateDNSNamespaceParameters `json:"forProvider"`
}

// PrivateDNSNamespaceObservation defines the observed state of PrivateDNSNamespace
type PrivateDNSNamespaceObservation struct {
	// A value that you can use to determine whether the request completed successfully.
	// To get the status of the operation, see GetOperation (https://docs.aws.amazon.com/cloud-map/latest/api/API_GetOperation.html).
	OperationID *string `json:"operationID,omitempty"`
}

// PrivateDNSNamespaceStatus defines the observed state of PrivateDNSNamespace.
type PrivateDNSNamespaceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PrivateDNSNamespaceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// PrivateDNSNamespace is the Schema for the PrivateDNSNamespaces API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type PrivateDNSNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PrivateDNSNamespaceSpec   `json:"spec"`
	Status            PrivateDNSNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PrivateDNSNamespaceList contains a list of PrivateDNSNamespaces
type PrivateDNSNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrivateDNSNamespace `json:"items"`
}

// Repository type metadata.
var (
	PrivateDNSNamespaceKind             = "PrivateDNSNamespace"
	PrivateDNSNamespaceGroupKind        = schema.GroupKind{Group: Group, Kind: PrivateDNSNamespaceKind}.String()
	PrivateDNSNamespaceKindAPIVersion   = PrivateDNSNamespaceKind + "." + GroupVersion.String()
	PrivateDNSNamespaceGroupVersionKind = GroupVersion.WithKind(PrivateDNSNamespaceKind)
)

func init() {
	SchemeBuilder.Register(&PrivateDNSNamespace{}, &PrivateDNSNamespaceList{})
}
