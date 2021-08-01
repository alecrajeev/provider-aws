/*
Copyright 2019 The Crossplane Authors.

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
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ec2 "github.com/crossplane/provider-aws/apis/ec2/v1beta1"
)

// ResolveReferences of this LoadBalancer
func (mg *LoadBalancer) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	// Resolve spec.forProvider.SecurityGroups
	mrsp, err := r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.SecurityGroups),
		References:    mg.Spec.ForProvider.SecurityGroupsRefs,
		Selector:      mg.Spec.ForProvider.SecurityGroupsSelector,
		To:            reference.To{Managed: &ec2.SecurityGroup{}, List: &ec2.SecurityGroupList{}},
		Extract:       reference.ExternalName(),
	})
	if err != nil {
		return errors.Wrap(err, "spec.forProvider.securityGroups")
	}
	mg.Spec.ForProvider.SecurityGroups = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.SecurityGroupsRefs = mrsp.ResolvedReferences

	// Resolve spec.forProvider.Subnets
	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.Subnets),
		References:    mg.Spec.ForProvider.SubnetsRefs,
		Selector:      mg.Spec.ForProvider.SubnetsSelector,
		To:            reference.To{Managed: &ec2.Subnet{}, List: &ec2.SubnetList{}},
		Extract:       reference.ExternalName(),
	})
	if err != nil {
		return errors.Wrap(err, "spec.forProvider.subnets")
	}
	mg.Spec.ForProvider.Subnets = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.SubnetsRefs = mrsp.ResolvedReferences

	// Resolve for Subnets and Elastic Ips in SubnetMappings
	for i := range mg.Spec.ForProvider.CustomSubnetMappingsParameters {
		rsp, err := r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].SubnetID),
			Reference:    mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].SubnetIDRefs,
			Selector:     mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].SubnetIDSelector,
			To:           reference.To{Managed: &ec2.Subnet{}, List: &ec2.SubnetList{}},
			Extract:      reference.ExternalName(),
		})
		if err != nil {
			return errors.Wrap(err, "spec.forProvider.subnetMappings.subnetID")
		}
		mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].SubnetID = reference.ToPtrValue(rsp.ResolvedValue)
		mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].SubnetIDRefs = rsp.ResolvedReference

		// Resolve for spec.forProvider.SubnetMappings[0].allocationID
		rspAddress, errAddress := r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].AllocationID),
			Reference:    mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].AllocationIDRefs,
			Selector:     mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].AllocationIDSelector,
			To:           reference.To{Managed: &ec2.Address{}, List: &ec2.AddressList{}},
			Extract:      reference.ExternalName(),
		})
		if errAddress != nil {
			return errors.Wrap(errAddress, "spec.forProvider.subnetMappings.allocationID")
		}
		mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].AllocationID = reference.ToPtrValue(rsp.ResolvedValue)
		mg.Spec.ForProvider.CustomSubnetMappingsParameters[i].AllocationIDRefs = rspAddress.ResolvedReference
	}

	return nil
}
