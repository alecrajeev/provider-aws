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

package loadbalancer

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"
	svcsdkapi "github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
	svcapitypes "github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
)

// SetupLoadBalancer adds a controller that reconciles LoadBalancer.
func SetupLoadBalancer(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter, poll time.Duration) error {
	name := managed.ControllerName(v1alpha1.LoadBalancerGroupKind)
	opts := []option{
		func(e *external) {
			e.postCreate = postCreate
			e.postObserve = postObserve
			e.preDelete = preDelete
			e.lateInitialize = lateInitialize
			e.isUpToDate = isUpToDate
			u := &updater{client: e.client}
			e.update = u.update
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.LoadBalancer{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.LoadBalancerGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), opts: opts}),
			managed.WithPollInterval(poll),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

// lateIntialize fills the empty fields in *svcapitypes.LoadBalancerParameters with
// the value seen in svcsdk.DescribeLoadBalancersOutput
func lateInitialize(cr *svcapitypes.LoadBalancerParameters, resp *svcsdk.DescribeLoadBalancersOutput) error {
	loadBalancer := resp.LoadBalancers[0]

	cr.IPAddressType = awsclient.LateInitializeStringPtr(cr.IPAddressType, loadBalancer.IpAddressType)
	cr.Scheme = awsclient.LateInitializeStringPtr(cr.Scheme, loadBalancer.Scheme)
	// for albs, the default security group gets set if not specified during creation
	if *cr.Type == "application" {
		cr.SecurityGroups = awsclient.LateInitializeStringPtrArray(cr.SecurityGroups, loadBalancer.SecurityGroups)
	}

	return nil
}

func postCreate(_ context.Context, cr *svcapitypes.LoadBalancer, obj *svcsdk.CreateLoadBalancerOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	// Set LoadBalancer's Arn as the external name annotation on the k8s object after creation.
	meta.SetExternalName(cr, aws.StringValue(obj.LoadBalancers[0].LoadBalancerArn))
	cre.ExternalNameAssigned = true
	return cre, nil
}

func postObserve(_ context.Context, cr *svcapitypes.LoadBalancer, resp *svcsdk.DescribeLoadBalancersOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	switch aws.StringValue(resp.LoadBalancers[0].State.Code) {
	case string(svcapitypes.LoadBalancerStateEnum_active):
		cr.SetConditions(xpv1.Available())
	case string(svcapitypes.LoadBalancerStateEnum_provisioning):
		cr.SetConditions(xpv1.Creating())
	case string(svcapitypes.LoadBalancerStateEnum_failed), string(svcapitypes.LoadBalancerStateEnum_active_impaired):
		cr.SetConditions(xpv1.Unavailable())
	}

	obs.ConnectionDetails = managed.ConnectionDetails{
		"name": []byte(aws.StringValue(resp.LoadBalancers[0].LoadBalancerArn)),
	}
	return obs, nil
}

func preDelete(_ context.Context, cr *svcapitypes.LoadBalancer, obj *svcsdk.DeleteLoadBalancerInput) (bool, error) {
	// Delete Load Balancer API call requires the ARN as a parameter,
	// so set the External Name to be an ARN.
	obj.LoadBalancerArn = aws.String(meta.GetExternalName(cr))
	return false, nil
}

func isUpToDate(cr *svcapitypes.LoadBalancer, obj *svcsdk.DescribeLoadBalancersOutput, objTags *svcsdk.DescribeTagsOutput) (bool, string, error) {

	diffIPAddressType := cmp.Diff(aws.StringValue(cr.Spec.ForProvider.IPAddressType), aws.StringValue(obj.LoadBalancers[0].IpAddressType))

	if diffIPAddressType != "" {
		return false, diffIPAddressType, nil
	}

	isUpToDateSecurityGroups, diffSecurityGroups := isUpToDateSecurityGroups(cr, obj)
	if !isUpToDateSecurityGroups {
		return false, diffSecurityGroups, nil
	}

	if len(cr.Spec.ForProvider.SubnetMappings) > 0 {
		isUpToDateSubnetMappings, diffSubnetMappings := isUpToDateSubnetMappings(cr, obj)
		if !isUpToDateSubnetMappings {
			return false, diffSubnetMappings, nil
		}
	} else {
		isUpToDateSubnets, diffSubnets := isUpToDateSubnets(cr, obj)
		if !isUpToDateSubnets {
			return false, diffSubnets, nil
		}
	}

	addTags, removeTags, diffTags := diffTags(cr.Spec.ForProvider.Tags, objTags)

	return len(addTags) == 0 && len(removeTags) == 0, diffTags, nil
}

func isUpToDateSecurityGroups(cr *svcapitypes.LoadBalancer, obj *svcsdk.DescribeLoadBalancersOutput) (bool, string) {
	// Handle nil pointer refs
	var securityGroups []*string
	var awsSecurityGroups []*string

	if cr.Spec.ForProvider.SecurityGroups != nil {
		securityGroups = cr.Spec.ForProvider.SecurityGroups
	}

	if obj.LoadBalancers[0].SecurityGroups != nil {
		awsSecurityGroups = obj.LoadBalancers[0].SecurityGroups
	}

	// Compare whether the slices are equal, ignore ordering
	sortCmp := cmpopts.SortSlices(func(i, j *string) bool {
		return aws.StringValue(i) < aws.StringValue(j)
	})

	diff := cmp.Diff(securityGroups, awsSecurityGroups, sortCmp, cmpopts.EquateEmpty())

	return diff == "", diff
}

func isUpToDateSubnets(cr *svcapitypes.LoadBalancer, obj *svcsdk.DescribeLoadBalancersOutput) (bool, string) {
	var loadBalancerAzs []*svcsdk.AvailabilityZone
	if len(obj.LoadBalancers) > 0 {
		loadBalancerAzs = obj.LoadBalancers[0].AvailabilityZones
	}

	// Handle nil pointer refs
	var subnets []*string
	awsSubnets := make([]*string, len(loadBalancerAzs))

	if cr.Spec.ForProvider.Subnets != nil {
		subnets = cr.Spec.ForProvider.Subnets
	}

	for i := range loadBalancerAzs {
		awsSubnets[i] = loadBalancerAzs[i].SubnetId
	}

	// Compare whether the slices are equal, ignore ordering
	sortCmp := cmpopts.SortSlices(func(i, j *string) bool {
		return aws.StringValue(i) < aws.StringValue(j)
	})

	diff := cmp.Diff(subnets, awsSubnets, sortCmp, cmpopts.EquateEmpty())

	return diff == "", diff
}

func isUpToDateSubnetMappings(cr *svcapitypes.LoadBalancer, obj *svcsdk.DescribeLoadBalancersOutput) (bool, string) {
	var loadBalancerAzs []*svcsdk.AvailabilityZone
	if len(obj.LoadBalancers) > 0 {
		loadBalancerAzs = obj.LoadBalancers[0].AvailabilityZones
	}

	// Handle nil pointer refs
	var subnetMappings []*svcapitypes.SubnetMapping
	awsSubnetMappings := make([]*svcapitypes.SubnetMapping, len(loadBalancerAzs))

	if cr.Spec.ForProvider.SubnetMappings != nil {
		subnetMappings = cr.Spec.ForProvider.SubnetMappings
	}

	for i := range loadBalancerAzs {
		// Define the SubnetMapping from the observed output.
		// Assumes only one LoadBalancer address because only
		// one is supported.
		var address *svcsdk.LoadBalancerAddress
		var allocationID *string
		var iPv6Address *string
		var privateIPv4Address *string
		if len(loadBalancerAzs[i].LoadBalancerAddresses) > 0 {
			address = loadBalancerAzs[i].LoadBalancerAddresses[0]
			allocationID = address.AllocationId
			iPv6Address = address.IPv6Address
			privateIPv4Address = address.PrivateIPv4Address
		}
		awsSubnetMappings[i] = &svcapitypes.SubnetMapping{
			AllocationID:       allocationID,
			IPv6Address:        iPv6Address,
			PrivateIPv4Address: privateIPv4Address,
			SubnetID:           loadBalancerAzs[i].SubnetId,
		}
	}

	// Compare whether the slices are equal, ignore ordering.
	// Sort by SubnetID alphabetically.
	sortCmp := cmpopts.SortSlices(func(s, p *svcapitypes.SubnetMapping) bool {
		return aws.StringValue(s.SubnetID) < aws.StringValue(p.SubnetID)
	})

	diff := cmp.Diff(subnetMappings, awsSubnetMappings, sortCmp, cmpopts.EquateEmpty())

	return diff == "", diff
}

// returns which AWS Tags exist in the resource tags and which are outdated and should be removed
func diffTags(spec []*svcapitypes.Tag, current *svcsdk.DescribeTagsOutput) (map[string]*string, []*string, string) {
	currentTags := GenerateMapFromTagsResponseOutput(current)
	specTags := GenerateMapFromTagsCR(spec)

	addMap := make(map[string]*string, len(specTags))
	removeTags := make([]*string, 0)

	for k, v := range currentTags {
		if awsclient.StringValue(specTags[k]) == awsclient.StringValue(v) {
			continue
		}
		removeTags = append(removeTags, awsclient.String(k))
	}
	for k, v := range specTags {
		if awsclient.StringValue(currentTags[k]) == awsclient.StringValue(v) {
			continue
		}
		addMap[k] = v
	}
	diffTags := ""
	if len(addMap) > 0 {
		diffTags += "AddTags: "
		for k, v := range addMap {
			diffTags += k + ": " + *v + ", "
		}
	}
	if len(removeTags) > 0 {
		diffTags += "\nRemoveTags: "
		for _, key := range removeTags {
			diffTags += *key + ", "
		}
	}

	return addMap, removeTags, diffTags
}

type updater struct {
	client svcsdkapi.ELBV2API
}

// nolint:gocyclo
func (u *updater) update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*svcapitypes.LoadBalancer)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#ELBV2.SetIpAddressType
	setIPAddressTypeInput := GenerateSetIPAddressTypeInput(cr)
	if _, err := u.client.SetIpAddressTypeWithContext(ctx, setIPAddressTypeInput); err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
	}

	// Only set security groups if they exist. Network Load Balancers don't have security groups.
	if len(cr.Spec.ForProvider.SecurityGroups) > 0 {
		// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#ELBV2.SetSecurityGroups
		setSecurityGroupsInput := GenerateSetSecurityGroupsInput(cr)
		if _, err := u.client.SetSecurityGroupsWithContext(ctx, setSecurityGroupsInput); err != nil {
			return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
		}
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#ELBV2.SetSubnets
	setSubnetsInput := GenerateSetSubnetsInput(cr)
	if _, err := u.client.SetSubnetsWithContext(ctx, setSubnetsInput); err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
	}

	// Tags
	tagInput := GenerateDescribeTagsInput(cr)
	tags, err := u.client.DescribeTagsWithContext(ctx, tagInput)
	if err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
	}

	addTags, removeTags, _ := diffTags(cr.Spec.ForProvider.Tags, tags)
	// Remove old tags before adding new tags in case values change for keys
	if len(removeTags) > 0 {
		tagsRemoveInput := GenerateRemoveTagsInput(removeTags, cr)
		if _, err := u.client.RemoveTagsWithContext(ctx, tagsRemoveInput); err != nil {
			return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
		}
	}
	if len(addTags) > 0 {
		tagsAddInput := GenerateAddTagsInput(addTags, cr)
		if _, err := u.client.AddTagsWithContext(ctx, tagsAddInput); err != nil {
			return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate)
		}
	}

	return managed.ExternalUpdate{}, nil
}

// GenerateSetIPAddressTypeInput is similar to GenerateCreateLoadBalancerInput
// Except it only sets the ip address type
func GenerateSetIPAddressTypeInput(cr *svcapitypes.LoadBalancer) *svcsdk.SetIpAddressTypeInput {
	f0 := &svcsdk.SetIpAddressTypeInput{}
	f0.SetLoadBalancerArn(meta.GetExternalName(cr))
	f0.SetIpAddressType(*cr.Spec.ForProvider.IPAddressType)

	return f0
}

// GenerateSetSecurityGroupsInput is similar to GenerateCreateLoadBalancerInput
// Except it only sets updated security groups
func GenerateSetSecurityGroupsInput(cr *svcapitypes.LoadBalancer) *svcsdk.SetSecurityGroupsInput {
	f0 := &svcsdk.SetSecurityGroupsInput{}
	f0.SetLoadBalancerArn(meta.GetExternalName(cr))
	f0.SetSecurityGroups(cr.Spec.ForProvider.SecurityGroups)

	return f0
}

// GenerateSetSubnetsInput is similar to GenerateCreaetLoadBalancerInput
// Except it sets the Subnets or SubnetMappings
func GenerateSetSubnetsInput(cr *svcapitypes.LoadBalancer) *svcsdk.SetSubnetsInput {
	f0 := &svcsdk.SetSubnetsInput{}
	f0.SetLoadBalancerArn(meta.GetExternalName(cr))
	if len(cr.Spec.ForProvider.SubnetMappings) > 0 {
		var subnetMappings []*svcsdk.SubnetMapping
		for _, subnetMapping := range cr.Spec.ForProvider.SubnetMappings {
			subnetMapping := svcsdk.SubnetMapping{
				AllocationId:       subnetMapping.AllocationID,
				IPv6Address:        subnetMapping.IPv6Address,
				PrivateIPv4Address: subnetMapping.PrivateIPv4Address,
				SubnetId:           subnetMapping.SubnetID,
			}
			subnetMappings = append(subnetMappings, &subnetMapping)
		}
		f0.SetSubnetMappings(subnetMappings)
	} else {
		f0.SetSubnets(cr.Spec.ForProvider.Subnets)
	}

	return f0
}
