package elasticsearchdomain

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	ctrl "sigs.k8s.io/controller-runtime"

	svcapitypes "github.com/crossplane/provider-aws/apis/elasticsearchservice/v1alpha1"
)

// SetupElasticsearchDomain adds a controller that reconciles ElasticsearchDomain.
func SetupElasticsearchDomain(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter, poll time.Duration) error {
	name := managed.ControllerName(svcapitypes.ElasticsearchDomainKind)
	opts := []option{
		func(e *external) {
			e.preObserve = preObserve
			e.postCreate = postCreate
			e.postObserve = postObserve
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewController(rl),
		}).
		For(&svcapitypes.ElasticsearchDomain{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(svcapitypes.ElasticsearchDomainGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), opts: opts}),
			managed.WithInitializers(),
			managed.WithPollInterval(poll),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

func preObserve(_ context.Context, cr *svcapitypes.ElasticsearchDomain, obj *svcsdk.DescribeElasticsearchDomainInput) error {
	obj.DomainName = aws.String(meta.GetExternalName(cr))
	return nil
}

// postCreate sets the external name annotation of the ElasticsearchDomain CRD after the create API call has been performed.
// The managed.WithInitializers() setting in the setup controller means that the external name annotation is not initially set.
func postCreate(_ context.Context, cr *svcapitypes.ElasticsearchDomain, resp *svcsdk.CreateElasticsearchDomainOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	meta.SetExternalName(cr, aws.StringValue(resp.DomainStatus.DomainName))
	return cre, nil
}

func postObserve(_ context.Context, cr *svcapitypes.ElasticsearchDomain, resp *svcsdk.DescribeElasticsearchDomainOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	if !*resp.DomainStatus.Created {
		cr.SetConditions(xpv1.Creating())
	} else {
		// Verify that both the Service Software and Elasticsearch version are not being upgraded.
		if !*resp.DomainStatus.UpgradeProcessing && *resp.DomainStatus.ServiceSoftwareOptions.UpdateStatus != string(svcapitypes.DeploymentStatus_IN_PROGRESS) {
			cr.SetConditions(xpv1.Available())
		} else {
			cr.SetConditions(xpv1.Unavailable())
		}
	}
	return obs, nil
}
