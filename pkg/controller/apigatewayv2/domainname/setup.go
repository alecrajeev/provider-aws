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

package domainname

import (
	"context"
	"time"

	svcsdk "github.com/aws/aws-sdk-go/service/apigatewayv2"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/event"
	"github.com/alecrajeev/crossplane-runtime/pkg/logging"
	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/ratelimiter"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"
	"github.com/alecrajeev/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane/provider-aws/apis/apigatewayv2/v1alpha1"
	aws "github.com/crossplane/provider-aws/pkg/clients"
)

// SetupDomainName adds a controller that reconciles DomainName.
func SetupDomainName(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter, poll time.Duration) error {
	name := managed.ControllerName(svcapitypes.DomainNameGroupKind)
	opts := []option{
		func(e *external) {
			e.preObserve = preObserve
			e.postObserve = postObserve
			e.preCreate = preCreate
			e.preDelete = preDelete
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&svcapitypes.DomainName{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(svcapitypes.DomainNameGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), opts: opts}),
			managed.WithPollInterval(poll),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

func preObserve(_ context.Context, cr *svcapitypes.DomainName, obj *svcsdk.GetDomainNameInput) error {
	obj.DomainName = aws.String(meta.GetExternalName(cr))
	return nil
}

func postObserve(_ context.Context, cr *svcapitypes.DomainName, _ *svcsdk.GetDomainNameOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	cr.SetConditions(xpv1.Available())
	return obs, nil
}

func preCreate(_ context.Context, cr *svcapitypes.DomainName, obj *svcsdk.CreateDomainNameInput) error {
	obj.DomainName = aws.String(meta.GetExternalName(cr))
	return nil
}

func preDelete(_ context.Context, cr *svcapitypes.DomainName, obj *svcsdk.DeleteDomainNameInput) (bool, error) {
	obj.DomainName = aws.String(meta.GetExternalName(cr))
	return false, nil
}
