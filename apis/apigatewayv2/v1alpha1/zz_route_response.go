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

// RouteResponseParameters defines the desired state of RouteResponse
type RouteResponseParameters struct {
	// Region is which region the RouteResponse will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	ModelSelectionExpression *string `json:"modelSelectionExpression,omitempty"`

	ResponseModels map[string]*string `json:"responseModels,omitempty"`

	ResponseParameters map[string]*ParameterConstraints `json:"responseParameters,omitempty"`

	// +kubebuilder:validation:Required
	RouteResponseKey              *string `json:"routeResponseKey"`
	CustomRouteResponseParameters `json:",inline"`
}

// RouteResponseSpec defines the desired state of RouteResponse
type RouteResponseSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RouteResponseParameters `json:"forProvider"`
}

// RouteResponseObservation defines the observed state of RouteResponse
type RouteResponseObservation struct {
	RouteResponseID *string `json:"routeResponseID,omitempty"`
}

// RouteResponseStatus defines the observed state of RouteResponse.
type RouteResponseStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RouteResponseObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// RouteResponse is the Schema for the RouteResponses API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type RouteResponse struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RouteResponseSpec   `json:"spec"`
	Status            RouteResponseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RouteResponseList contains a list of RouteResponses
type RouteResponseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RouteResponse `json:"items"`
}

// Repository type metadata.
var (
	RouteResponseKind             = "RouteResponse"
	RouteResponseGroupKind        = schema.GroupKind{Group: Group, Kind: RouteResponseKind}.String()
	RouteResponseKindAPIVersion   = RouteResponseKind + "." + GroupVersion.String()
	RouteResponseGroupVersionKind = GroupVersion.WithKind(RouteResponseKind)
)

func init() {
	SchemeBuilder.Register(&RouteResponse{}, &RouteResponseList{})
}
