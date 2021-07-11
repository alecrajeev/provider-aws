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

package redshift

import (
	"context"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsredshift "github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/event"
	"github.com/alecrajeev/crossplane-runtime/pkg/logging"
	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/password"
	"github.com/alecrajeev/crossplane-runtime/pkg/ratelimiter"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"
	"github.com/alecrajeev/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-aws/apis/redshift/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/redshift"
)

const (
	errUnexpectedObject = "managed resource is not a Redshift custom resource"
	errKubeUpdateFailed = "cannot update Redshift cluster custom resource"
	errMultipleCluster  = "multiple clusters with the same name found"
	errCreateFailed     = "cannot create Redshift cluster"
	errModifyFailed     = "cannot modify Redshift cluster"
	errDeleteFailed     = "cannot delete Redshift cluster"
	errDescribeFailed   = "cannot describe Redshift cluster"
	errUpToDateFailed   = "cannot check whether object is up-to-date"
)

// SetupCluster adds a controller that reconciles Redshift clusters.
func SetupCluster(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter, poll time.Duration) error {
	name := managed.ControllerName(v1alpha1.ClusterGroupKind)
	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Cluster{}).
		Complete(managed.NewReconciler(
			mgr, resource.ManagedKind(v1alpha1.ClusterGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: redshift.NewClient}),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithPollInterval(poll),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type connector struct {
	kube        client.Client
	newClientFn func(config aws.Config) redshift.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Cluster)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := awsclient.GetConfig(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, err
	}
	return &external{client: c.newClientFn(*cfg), kube: c.kube}, nil
}

type external struct {
	kube   client.Client
	client redshift.Client
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) { //nolint:gocyclo
	cr, ok := mg.(*v1alpha1.Cluster)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	rsp, err := e.client.DescribeClustersRequest(&awsredshift.DescribeClustersInput{
		ClusterIdentifier: aws.String(meta.GetExternalName(cr)),
	}).Send(ctx)
	if err != nil {
		return managed.ExternalObservation{}, awsclient.Wrap(resource.Ignore(redshift.IsNotFound, err), errDescribeFailed)
	}

	// Describe requests can be used with filters, which then returns a list.
	// But we use an explicit identifier, so, if there is no error, there should
	// be only 1 element in the list.
	if len(rsp.Clusters) != 1 {
		return managed.ExternalObservation{}, errors.New(errMultipleCluster)
	}
	instance := rsp.Clusters[0]
	current := cr.Spec.ForProvider.DeepCopy()
	redshift.LateInitialize(&cr.Spec.ForProvider, &instance)
	if !reflect.DeepEqual(current, &cr.Spec.ForProvider) {
		if err := e.kube.Update(ctx, cr); err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errKubeUpdateFailed)
		}
	}

	cr.Status.AtProvider = redshift.GenerateObservation(rsp.Clusters[0])
	switch cr.Status.AtProvider.ClusterStatus {
	case v1alpha1.StateAvailable:
		cr.Status.SetConditions(xpv1.Available())
	case v1alpha1.StateCreating:
		cr.Status.SetConditions(xpv1.Creating())
	case v1alpha1.StateDeleting:
		cr.Status.SetConditions(xpv1.Deleting())
	default:
		cr.Status.SetConditions(xpv1.Unavailable())
	}

	updated, err := redshift.IsUpToDate(cr.Spec.ForProvider, instance)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errUpToDateFailed)
	}

	return managed.ExternalObservation{
		ResourceUpToDate:  updated,
		ResourceExists:    true,
		ConnectionDetails: redshift.GetConnectionDetails(*cr),
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Cluster)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.SetConditions(xpv1.Creating())
	if cr.Status.AtProvider.ClusterStatus == v1alpha1.StateCreating {
		return managed.ExternalCreation{}, nil
	}
	pw, err := password.Generate()
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	input := redshift.GenerateCreateClusterInput(&cr.Spec.ForProvider, aws.String(meta.GetExternalName(cr)), aws.String(pw))
	_, err = e.client.CreateClusterRequest(input).Send(ctx)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreateFailed)
	}

	conn := managed.ConnectionDetails{
		xpv1.ResourceCredentialsSecretPasswordKey: []byte(aws.StringValue(input.MasterUserPassword)),
		xpv1.ResourceCredentialsSecretUserKey:     []byte(aws.StringValue(input.MasterUsername)),
	}

	return managed.ExternalCreation{ConnectionDetails: conn}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Cluster)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}
	switch cr.Status.AtProvider.ClusterStatus {
	case v1alpha1.StateModifying, v1alpha1.StateCreating:
		return managed.ExternalUpdate{}, nil
	}

	rsp, err := e.client.DescribeClustersRequest(&awsredshift.DescribeClustersInput{
		ClusterIdentifier: aws.String(meta.GetExternalName(cr)),
	}).Send(ctx)
	if err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(resource.Ignore(redshift.IsNotFound, err), errDescribeFailed)
	}

	_, err = e.client.ModifyClusterRequest(redshift.GenerateModifyClusterInput(&cr.Spec.ForProvider, rsp.Clusters[0])).Send(ctx)

	if err == nil && aws.StringValue(cr.Spec.ForProvider.NewClusterIdentifier) != meta.GetExternalName(cr) {
		meta.SetExternalName(cr, aws.StringValue(cr.Spec.ForProvider.NewClusterIdentifier))

		if err := e.kube.Update(ctx, cr); err != nil {
			return managed.ExternalUpdate{}, errors.Wrap(err, errKubeUpdateFailed)
		}
	}

	return managed.ExternalUpdate{}, awsclient.Wrap(err, errModifyFailed)
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Cluster)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.SetConditions(xpv1.Deleting())
	if cr.Status.AtProvider.ClusterStatus == v1alpha1.StateDeleting {
		return nil
	}

	_, err := e.client.DeleteClusterRequest(redshift.GenerateDeleteClusterInput(&cr.Spec.ForProvider, aws.String(meta.GetExternalName(cr)))).Send(ctx)

	return awsclient.Wrap(resource.Ignore(redshift.IsNotFound, err), errDeleteFailed)
}
