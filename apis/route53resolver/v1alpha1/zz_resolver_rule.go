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

// ResolverRuleParameters defines the desired state of ResolverRule
type ResolverRuleParameters struct {
	// Region is which region the ResolverRule will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// DNS queries for this domain name are forwarded to the IP addresses that you
	// specify in TargetIps. If a query matches multiple Resolver rules (example.com
	// and www.example.com), outbound DNS queries are routed using the Resolver
	// rule that contains the most specific domain name (www.example.com).
	// +kubebuilder:validation:Required
	DomainName *string `json:"domainName"`
	// A friendly name that lets you easily find a rule in the Resolver dashboard
	// in the Route 53 console.
	Name *string `json:"name,omitempty"`
	// The ID of the outbound Resolver endpoint that you want to use to route DNS
	// queries to the IP addresses that you specify in TargetIps.
	ResolverEndpointID *string `json:"resolverEndpointID,omitempty"`
	// When you want to forward DNS queries for specified domain name to resolvers
	// on your network, specify FORWARD.
	//
	// When you have a forwarding rule to forward DNS queries for a domain to your
	// network and you want Resolver to process queries for a subdomain of that
	// domain, specify SYSTEM.
	//
	// For example, to forward DNS queries for example.com to resolvers on your
	// network, you create a rule and specify FORWARD for RuleType. To then have
	// Resolver process queries for apex.example.com, you create a rule and specify
	// SYSTEM for RuleType.
	//
	// Currently, only Resolver can create rules that have a value of RECURSIVE
	// for RuleType.
	// +kubebuilder:validation:Required
	RuleType *string `json:"ruleType"`
	// A list of the tag keys and values that you want to associate with the endpoint.
	Tags []*Tag `json:"tags,omitempty"`
	// The IPs that you want Resolver to forward DNS queries to. You can specify
	// only IPv4 addresses. Separate IP addresses with a comma.
	//
	// TargetIps is available only when the value of Rule type is FORWARD.
	TargetIPs                    []*TargetAddress `json:"targetIPs,omitempty"`
	CustomResolverRuleParameters `json:",inline"`
}

// ResolverRuleSpec defines the desired state of ResolverRule
type ResolverRuleSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ResolverRuleParameters `json:"forProvider"`
}

// ResolverRuleObservation defines the observed state of ResolverRule
type ResolverRuleObservation struct {
	// The ARN (Amazon Resource Name) for the Resolver rule specified by Id.
	ARN *string `json:"arn,omitempty"`
	// The date and time that the Resolver rule was created, in Unix time format
	// and Coordinated Universal Time (UTC).
	CreationTime *string `json:"creationTime,omitempty"`
	// A unique string that you specified when you created the Resolver rule. CreatorRequestId
	// identifies the request and allows failed requests to be retried without the
	// risk of executing the operation twice.
	CreatorRequestID *string `json:"creatorRequestID,omitempty"`
	// The ID that Resolver assigned to the Resolver rule when you created it.
	ID *string `json:"id,omitempty"`
	// The date and time that the Resolver rule was last updated, in Unix time format
	// and Coordinated Universal Time (UTC).
	ModificationTime *string `json:"modificationTime,omitempty"`
	// When a rule is shared with another AWS account, the account ID of the account
	// that the rule is shared with.
	OwnerID *string `json:"ownerID,omitempty"`
	// Whether the rules is shared and, if so, whether the current account is sharing
	// the rule with another account, or another account is sharing the rule with
	// the current account.
	ShareStatus *string `json:"shareStatus,omitempty"`
	// A code that specifies the current status of the Resolver rule.
	Status *string `json:"status,omitempty"`
	// A detailed description of the status of a Resolver rule.
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// ResolverRuleStatus defines the observed state of ResolverRule.
type ResolverRuleStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ResolverRuleObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// ResolverRule is the Schema for the ResolverRules API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type ResolverRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ResolverRuleSpec   `json:"spec,omitempty"`
	Status            ResolverRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ResolverRuleList contains a list of ResolverRules
type ResolverRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResolverRule `json:"items"`
}

// Repository type metadata.
var (
	ResolverRuleKind             = "ResolverRule"
	ResolverRuleGroupKind        = schema.GroupKind{Group: Group, Kind: ResolverRuleKind}.String()
	ResolverRuleKindAPIVersion   = ResolverRuleKind + "." + GroupVersion.String()
	ResolverRuleGroupVersionKind = GroupVersion.WithKind(ResolverRuleKind)
)

func init() {
	SchemeBuilder.Register(&ResolverRule{}, &ResolverRuleList{})
}
