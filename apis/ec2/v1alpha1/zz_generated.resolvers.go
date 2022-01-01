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

// Code generated by angryjet. DO NOT EDIT.

package v1alpha1

import (
	"context"

	reference "github.com/crossplane/crossplane-runtime/pkg/reference"
	manualv1alpha1 "github.com/crossplane/provider-aws/apis/ec2/manualv1alpha1"
	v1beta1 "github.com/crossplane/provider-aws/apis/ec2/v1beta1"
	v1alpha1 "github.com/crossplane/provider-aws/apis/kms/v1alpha1"
	errors "github.com/pkg/errors"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// ResolveReferences of this LaunchTemplateVersion.
func (mg *LaunchTemplateVersion) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateIDRef,
		Selector:     mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateIDSelector,
		To: reference.To{
			List:    &LaunchTemplateList{},
			Managed: &LaunchTemplate{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateID")
	}
	mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateName),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateNameRef,
		Selector:     mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateNameSelector,
		To: reference.To{
			List:    &LaunchTemplateList{},
			Managed: &LaunchTemplate{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateName")
	}
	mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateName = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomLaunchTemplateVersionParameters.LaunchTemplateNameRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this Route.
func (mg *Route) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayIDRef,
		Selector:     mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayIDSelector,
		To: reference.To{
			List:    &TransitGatewayList{},
			Managed: &TransitGateway{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayID")
	}
	mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomRouteParameters.TransitGatewayIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomRouteParameters.NATGatewayID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomRouteParameters.NATGatewayIDRef,
		Selector:     mg.Spec.ForProvider.CustomRouteParameters.NATGatewayIDSelector,
		To: reference.To{
			List:    &v1beta1.NATGatewayList{},
			Managed: &v1beta1.NATGateway{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomRouteParameters.NATGatewayID")
	}
	mg.Spec.ForProvider.CustomRouteParameters.NATGatewayID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomRouteParameters.NATGatewayIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionIDRef,
		Selector:     mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionIDSelector,
		To: reference.To{
			List:    &VPCPeeringConnectionList{},
			Managed: &VPCPeeringConnection{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionID")
	}
	mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomRouteParameters.VPCPeeringConnectionIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomRouteParameters.InstanceID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomRouteParameters.InstanceIDRef,
		Selector:     mg.Spec.ForProvider.CustomRouteParameters.InstanceIDSelector,
		To: reference.To{
			List:    &manualv1alpha1.InstanceList{},
			Managed: &manualv1alpha1.Instance{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomRouteParameters.InstanceID")
	}
	mg.Spec.ForProvider.CustomRouteParameters.InstanceID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomRouteParameters.InstanceIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomRouteParameters.GatewayID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomRouteParameters.GatewayIDRef,
		Selector:     mg.Spec.ForProvider.CustomRouteParameters.GatewayIDSelector,
		To: reference.To{
			List:    &v1beta1.InternetGatewayList{},
			Managed: &v1beta1.InternetGateway{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomRouteParameters.GatewayID")
	}
	mg.Spec.ForProvider.CustomRouteParameters.GatewayID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomRouteParameters.GatewayIDRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this TransitGatewayRoute.
func (mg *TransitGatewayRoute) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentIDRef,
		Selector:     mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentIDSelector,
		To: reference.To{
			List:    &TransitGatewayVPCAttachmentList{},
			Managed: &TransitGatewayVPCAttachment{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentID")
	}
	mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayAttachmentIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableIDRef,
		Selector:     mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableIDSelector,
		To: reference.To{
			List:    &TransitGatewayRouteTableList{},
			Managed: &TransitGatewayRouteTable{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableID")
	}
	mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomTransitGatewayRouteParameters.TransitGatewayRouteTableIDRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this TransitGatewayRouteTable.
func (mg *TransitGatewayRouteTable) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayIDRef,
		Selector:     mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayIDSelector,
		To: reference.To{
			List:    &TransitGatewayList{},
			Managed: &TransitGateway{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayID")
	}
	mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomTransitGatewayRouteTableParameters.TransitGatewayIDRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this TransitGatewayVPCAttachment.
func (mg *TransitGatewayVPCAttachment) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var mrsp reference.MultiResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCIDRef,
		Selector:     mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCIDSelector,
		To: reference.To{
			List:    &v1beta1.VPCList{},
			Managed: &v1beta1.VPC{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCID")
	}
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.VPCIDRef = rsp.ResolvedReference

	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDs),
		Extract:       reference.ExternalName(),
		References:    mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDRefs,
		Selector:      mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDSelector,
		To: reference.To{
			List:    &v1beta1.SubnetList{},
			Managed: &v1beta1.Subnet{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDs")
	}
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDs = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.SubnetIDRefs = mrsp.ResolvedReferences

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayIDRef,
		Selector:     mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayIDSelector,
		To: reference.To{
			List:    &TransitGatewayList{},
			Managed: &TransitGateway{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayID")
	}
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomTransitGatewayVPCAttachmentParameters.TransitGatewayIDRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this VPCEndpoint.
func (mg *VPCEndpoint) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var mrsp reference.MultiResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCIDRef,
		Selector:     mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCIDSelector,
		To: reference.To{
			List:    &v1beta1.VPCList{},
			Managed: &v1beta1.VPC{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCID")
	}
	mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomVPCEndpointParameters.VPCIDRef = rsp.ResolvedReference

	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDs),
		Extract:       reference.ExternalName(),
		References:    mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDRefs,
		Selector:      mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDSelector,
		To: reference.To{
			List:    &v1beta1.SecurityGroupList{},
			Managed: &v1beta1.SecurityGroup{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDs")
	}
	mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDs = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.CustomVPCEndpointParameters.SecurityGroupIDRefs = mrsp.ResolvedReferences

	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDs),
		Extract:       reference.ExternalName(),
		References:    mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDRefs,
		Selector:      mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDSelector,
		To: reference.To{
			List:    &v1beta1.SubnetList{},
			Managed: &v1beta1.Subnet{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDs")
	}
	mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDs = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.CustomVPCEndpointParameters.SubnetIDRefs = mrsp.ResolvedReferences

	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDs),
		Extract:       reference.ExternalName(),
		References:    mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDRefs,
		Selector:      mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDSelector,
		To: reference.To{
			List:    &v1beta1.RouteTableList{},
			Managed: &v1beta1.RouteTable{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDs")
	}
	mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDs = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.CustomVPCEndpointParameters.RouteTableIDRefs = mrsp.ResolvedReferences

	return nil
}

// ResolveReferences of this VPCPeeringConnection.
func (mg *VPCPeeringConnection) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCIDRef,
		Selector:     mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCIDSelector,
		To: reference.To{
			List:    &v1beta1.VPCList{},
			Managed: &v1beta1.VPC{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCID")
	}
	mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomVPCPeeringConnectionParameters.VPCIDRef = rsp.ResolvedReference

	return nil
}

// ResolveReferences of this Volume.
func (mg *Volume) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyID),
		Extract:      reference.ExternalName(),
		Reference:    mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyIDRef,
		Selector:     mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyIDSelector,
		To: reference.To{
			List:    &v1alpha1.KeyList{},
			Managed: &v1alpha1.Key{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyID")
	}
	mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomVolumeParameters.KMSKeyIDRef = rsp.ResolvedReference

	return nil
}
