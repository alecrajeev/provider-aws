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

// Code originally genrated by ack-generate.
// Edited to support tags and diff messages.

package loadbalancer

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"
	svcsdkapi "github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
	svcapitypes "github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an LoadBalancer resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create LoadBalancer in AWS"
	errUpdate        = "cannot update LoadBalancer in AWS"
	errDescribe      = "failed to describe LoadBalancer"
	errDescribeTags  = "failed to describe LoadBalancer tags"
	errDelete        = "failed to delete LoadBalancer"
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
			cpresource.ManagedKind(v1alpha1.LoadBalancerGroupVersionKind),
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
func (u *updater) update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
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

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.LoadBalancer)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := awsclient.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcsdk.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.LoadBalancer)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeLoadBalancersInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeLoadBalancersWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}

	// Only able to get tags if there is a valid arn and a response
	// with an existing load balancer. DescribeTagsWithContext requires
	// a valid arn and the Load Balancer resource to exist to be
	// successful.
	var respTags *svcsdk.DescribeTagsOutput
	var errTags error
	validArn := false
	if strings.Contains(meta.GetExternalName(cr), "arn") && len(resp.LoadBalancers) > 0 {
		tagInput := GenerateDescribeTagsInput(cr)
		respTags, errTags = e.client.DescribeTagsWithContext(ctx, tagInput)
		if errTags != nil {
			return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, errTags), errDescribeTags)
		}
		validArn = true
	}

	resp = e.filterList(cr, resp)
	if len(resp.LoadBalancers) == 0 {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}

	// Only copy tags into custom resource if there is a valid arn.
	if validArn {
		GenerateLoadBalancerWithTags(resp, respTags).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	} else {
		GenerateLoadBalancer(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	}

	upToDate, diff, err := e.isUpToDate(cr, resp, respTags)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
		Diff:                    diff,
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.LoadBalancer)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateLoadBalancerInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateLoadBalancerWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreate)
	}

	if resp.LoadBalancers != nil {
		f0 := []*svcapitypes.LoadBalancer_SDK{}
		for _, f0iter := range resp.LoadBalancers {
			f0elem := &svcapitypes.LoadBalancer_SDK{}
			if f0iter.AvailabilityZones != nil {
				f0elemf0 := []*svcapitypes.AvailabilityZone{}
				for _, f0elemf0iter := range f0iter.AvailabilityZones {
					f0elemf0elem := &svcapitypes.AvailabilityZone{}
					if f0elemf0iter.LoadBalancerAddresses != nil {
						f0elemf0elemf0 := []*svcapitypes.LoadBalancerAddress{}
						for _, f0elemf0elemf0iter := range f0elemf0iter.LoadBalancerAddresses {
							f0elemf0elemf0elem := &svcapitypes.LoadBalancerAddress{}
							if f0elemf0elemf0iter.AllocationId != nil {
								f0elemf0elemf0elem.AllocationID = f0elemf0elemf0iter.AllocationId
							}
							if f0elemf0elemf0iter.IPv6Address != nil {
								f0elemf0elemf0elem.IPv6Address = f0elemf0elemf0iter.IPv6Address
							}
							if f0elemf0elemf0iter.IpAddress != nil {
								f0elemf0elemf0elem.IPAddress = f0elemf0elemf0iter.IpAddress
							}
							if f0elemf0elemf0iter.PrivateIPv4Address != nil {
								f0elemf0elemf0elem.PrivateIPv4Address = f0elemf0elemf0iter.PrivateIPv4Address
							}
							f0elemf0elemf0 = append(f0elemf0elemf0, f0elemf0elemf0elem)
						}
						f0elemf0elem.LoadBalancerAddresses = f0elemf0elemf0
					}
					if f0elemf0iter.OutpostId != nil {
						f0elemf0elem.OutpostID = f0elemf0iter.OutpostId
					}
					if f0elemf0iter.SubnetId != nil {
						f0elemf0elem.SubnetID = f0elemf0iter.SubnetId
					}
					if f0elemf0iter.ZoneName != nil {
						f0elemf0elem.ZoneName = f0elemf0iter.ZoneName
					}
					f0elemf0 = append(f0elemf0, f0elemf0elem)
				}
				f0elem.AvailabilityZones = f0elemf0
			}
			if f0iter.CanonicalHostedZoneId != nil {
				f0elem.CanonicalHostedZoneID = f0iter.CanonicalHostedZoneId
			}
			if f0iter.CreatedTime != nil {
				f0elem.CreatedTime = &metav1.Time{*f0iter.CreatedTime}
			}
			if f0iter.CustomerOwnedIpv4Pool != nil {
				f0elem.CustomerOwnedIPv4Pool = f0iter.CustomerOwnedIpv4Pool
			}
			if f0iter.DNSName != nil {
				f0elem.DNSName = f0iter.DNSName
			}
			if f0iter.IpAddressType != nil {
				f0elem.IPAddressType = f0iter.IpAddressType
			}
			if f0iter.LoadBalancerArn != nil {
				f0elem.LoadBalancerARN = f0iter.LoadBalancerArn
			}
			if f0iter.LoadBalancerName != nil {
				f0elem.LoadBalancerName = f0iter.LoadBalancerName
			}
			if f0iter.Scheme != nil {
				f0elem.Scheme = f0iter.Scheme
			}
			if f0iter.SecurityGroups != nil {
				f0elemf9 := []*string{}
				for _, f0elemf9iter := range f0iter.SecurityGroups {
					var f0elemf9elem string
					f0elemf9elem = *f0elemf9iter
					f0elemf9 = append(f0elemf9, &f0elemf9elem)
				}
				f0elem.SecurityGroups = f0elemf9
			}
			if f0iter.State != nil {
				f0elemf10 := &svcapitypes.LoadBalancerState{}
				if f0iter.State.Code != nil {
					f0elemf10.Code = f0iter.State.Code
				}
				if f0iter.State.Reason != nil {
					f0elemf10.Reason = f0iter.State.Reason
				}
				f0elem.State = f0elemf10
			}
			if f0iter.Type != nil {
				f0elem.Type = f0iter.Type
			}
			if f0iter.VpcId != nil {
				f0elem.VPCID = f0iter.VpcId
			}
			f0 = append(f0, f0elem)
		}
		cr.Status.AtProvider.LoadBalancers = f0
	} else {
		cr.Status.AtProvider.LoadBalancers = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	return e.update(ctx, mg)

}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.LoadBalancer)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteLoadBalancerInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteLoadBalancerWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.ELBV2API, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		filterList:     nopFilterList,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		update:         nopUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.ELBV2API
	preObserve     func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersInput) error
	postObserve    func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	filterList     func(*svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersOutput) *svcsdk.DescribeLoadBalancersOutput
	lateInitialize func(*svcapitypes.LoadBalancerParameters, *svcsdk.DescribeLoadBalancersOutput) error
	isUpToDate     func(*svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersOutput, *svcsdk.DescribeTagsOutput) (bool, string, error)
	preCreate      func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.CreateLoadBalancerInput) error
	postCreate     func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.CreateLoadBalancerOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DeleteLoadBalancerInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DeleteLoadBalancerOutput, error) error
	update         func(context.Context, cpresource.Managed) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersInput) error {
	return nil
}
func nopPostObserve(_ context.Context, _ *svcapitypes.LoadBalancer, _ *svcsdk.DescribeLoadBalancersOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopFilterList(_ *svcapitypes.LoadBalancer, list *svcsdk.DescribeLoadBalancersOutput) *svcsdk.DescribeLoadBalancersOutput {
	return list
}

func nopLateInitialize(*svcapitypes.LoadBalancerParameters, *svcsdk.DescribeLoadBalancersOutput) error {
	return nil
}
func alwaysUpToDate(*svcapitypes.LoadBalancer, *svcsdk.DescribeLoadBalancersOutput, *svcsdk.DescribeTagsOutput) (bool, string, error) {
	return true, "", nil
}

func nopPreCreate(context.Context, *svcapitypes.LoadBalancer, *svcsdk.CreateLoadBalancerInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.LoadBalancer, _ *svcsdk.CreateLoadBalancerOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.LoadBalancer, *svcsdk.DeleteLoadBalancerInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.LoadBalancer, _ *svcsdk.DeleteLoadBalancerOutput, err error) error {
	return err
}
func nopUpdate(context.Context, cpresource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}
