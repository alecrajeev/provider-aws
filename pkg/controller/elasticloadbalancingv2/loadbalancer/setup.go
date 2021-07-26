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

/*
TODO: Fix Order

	<stdlib pkgs>

	<external pkgs (anything not in github.com/crossplane)>

	<crossplane org pkgs>

	<local to this repo pkgs>

*/

import (
	"context"
	"time"

	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	awsclient "github.com/crossplane/provider-aws/pkg/clients"

	"github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
	svcapitypes "github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
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
