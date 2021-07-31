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
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// RuleParameters defines the desired state of Rule
type RuleParameters struct {
	// Region is which region the Rule will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// The actions.
	// +kubebuilder:validation:Required
	Actions []*Action `json:"actions"`
	// The conditions.
	// +kubebuilder:validation:Required
	Conditions []*RuleCondition `json:"conditions"`
	// The Amazon Resource Name (ARN) of the listener.
	// +kubebuilder:validation:Required
	ListenerARN *string `json:"listenerARN"`
	// The rule priority. A listener can't have multiple rules with the same priority.
	// +kubebuilder:validation:Required
	Priority *int64 `json:"priority"`
	// The tags to assign to the rule.
	Tags                 []*Tag `json:"tags,omitempty"`
	CustomRuleParameters `json:",inline"`
}

// RuleSpec defines the desired state of Rule
type RuleSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RuleParameters `json:"forProvider"`
}

// RuleObservation defines the observed state of Rule
type RuleObservation struct {
	// Information about the rule.
	Rules []*Rule_SDK `json:"rules,omitempty"`
}

// RuleStatus defines the observed state of Rule.
type RuleStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RuleObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Rule is the Schema for the Rules API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type Rule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RuleSpec   `json:"spec"`
	Status            RuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RuleList contains a list of Rules
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rule `json:"items"`
}

// Repository type metadata.
var (
	RuleKind             = "Rule"
	RuleGroupKind        = schema.GroupKind{Group: Group, Kind: RuleKind}.String()
	RuleKindAPIVersion   = RuleKind + "." + GroupVersion.String()
	RuleGroupVersionKind = GroupVersion.WithKind(RuleKind)
)

func init() {
	SchemeBuilder.Register(&Rule{}, &RuleList{})
}
