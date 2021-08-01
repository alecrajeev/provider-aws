/*
Copyright 2020 The Crossplane Authors.

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

import xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

// CustomLoadBalancerParameters includes custom fields for LoadBalancerParameters.
type CustomLoadBalancerParameters struct {
	SecurityGroups []*string `json:"securityGroups,omitempty"`

	// SecurityGroupsRef is a list of references to SecurityGroups used to set
	// the SecurityGroups.
	// +optional
	SecurityGroupsRefs []xpv1.Reference `json:"securityGroupsRefs,omitempty"`

	// SecurityGroupSelector selects references to SecurityGroups
	// +optional
	SecurityGroupsSelector *xpv1.Selector `json:"securityGroupsSelector,omitempty"`

	Subnets []*string `json:"subnets,omitempty"`

	// SubnetsRef is a list of references to Subnets
	// +optional
	SubnetsRefs []xpv1.Reference `json:"subnetsRefs,omitempty"`

	// SubnetsSelector selects references to Subnets
	// +optional
	SubnetsSelector *xpv1.Selector `json:"subnetSelector,omitempty"`

	// CustomSubnetMappingsParameters are subnet mappings
	// +optional
	CustomSubnetMappingsParameters []CustomSubnetMappingParameter `json:"subnetMappings,omitempty"`
}

// CustomSubnetMappingParameter includes custom fields for a SubnetMapping
type CustomSubnetMappingParameter struct {
	IPv6Address *string `json:"iPv6Address,omitempty"`

	PrivateIPv4Address *string `json:"privateIPv4Address,omitempty"`

	SubnetID *string `json:"subnetID,omitempty"`

	// SubnetIDRefs is a references to a Subnet.
	// +optional
	SubnetIDRefs *xpv1.Reference `json:"subnetIDRefs,omitempty"`

	// SubnetIDSelector selects references to a Subnet.
	// +optional
	SubnetIDSelector *xpv1.Selector `json:"subnetIDSelector,omitempty"`

	AllocationID *string `json:"allocationID,omitempty"`

	// AllocationIDRefs is a reference to an AllocationID
	// that corresponds to an elastic ip address.
	// +optional
	AllocationIDRefs *xpv1.Reference `json:"allocationIDRefs,omitempty"`

	// AllocationIDSelector is a selector to an AllocationID
	// that corresponds to an elastic ip address.
	// +optional
	AllocationIDSelector *xpv1.Selector `json:"allocationIDSelector,omitempty"`
}
