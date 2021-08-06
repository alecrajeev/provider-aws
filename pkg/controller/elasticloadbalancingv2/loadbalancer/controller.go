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
// Edited Observe to support tags and diff messages.

package loadbalancer

import (
	"context"
	"strings"

	svcapi "github.com/aws/aws-sdk-go/service/elbv2"
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"
	svcsdkapi "github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

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
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
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
