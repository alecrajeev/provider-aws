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

// Code generated by ack-generate. DO NOT EDIT.

package secret

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/secretsmanager"
	svcsdk "github.com/aws/aws-sdk-go/service/secretsmanager"
	svcsdkapi "github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/alecrajeev/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane/provider-aws/apis/secretsmanager/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an Secret resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create Secret in AWS"
	errUpdate        = "cannot update Secret in AWS"
	errDescribe      = "failed to describe Secret"
	errDelete        = "failed to delete Secret"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.Secret)
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
	cr, ok := mg.(*svcapitypes.Secret)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeSecretInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeSecretWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateSecret(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	upToDate, err := e.isUpToDate(cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.Secret)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateSecretInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateSecretWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreate)
	}

	if resp.ARN != nil {
		cr.Status.AtProvider.ARN = resp.ARN
	} else {
		cr.Status.AtProvider.ARN = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*svcapitypes.Secret)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}
	input := GenerateUpdateSecretInput(cr)
	if err := e.preUpdate(ctx, cr, input); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "pre-update failed")
	}
	resp, err := e.client.UpdateSecretWithContext(ctx, input)
	return e.postUpdate(ctx, cr, resp, managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdate))
}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.Secret)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteSecretInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteSecretWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.SecretsManagerAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		preUpdate:      nopPreUpdate,
		postUpdate:     nopPostUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.SecretsManagerAPI
	preObserve     func(context.Context, *svcapitypes.Secret, *svcsdk.DescribeSecretInput) error
	postObserve    func(context.Context, *svcapitypes.Secret, *svcsdk.DescribeSecretOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.SecretParameters, *svcsdk.DescribeSecretOutput) error
	isUpToDate     func(*svcapitypes.Secret, *svcsdk.DescribeSecretOutput) (bool, error)
	preCreate      func(context.Context, *svcapitypes.Secret, *svcsdk.CreateSecretInput) error
	postCreate     func(context.Context, *svcapitypes.Secret, *svcsdk.CreateSecretOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.Secret, *svcsdk.DeleteSecretInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.Secret, *svcsdk.DeleteSecretOutput, error) error
	preUpdate      func(context.Context, *svcapitypes.Secret, *svcsdk.UpdateSecretInput) error
	postUpdate     func(context.Context, *svcapitypes.Secret, *svcsdk.UpdateSecretOutput, managed.ExternalUpdate, error) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.Secret, *svcsdk.DescribeSecretInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.Secret, _ *svcsdk.DescribeSecretOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.SecretParameters, *svcsdk.DescribeSecretOutput) error {
	return nil
}
func alwaysUpToDate(*svcapitypes.Secret, *svcsdk.DescribeSecretOutput) (bool, error) {
	return true, nil
}

func nopPreCreate(context.Context, *svcapitypes.Secret, *svcsdk.CreateSecretInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.Secret, _ *svcsdk.CreateSecretOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.Secret, *svcsdk.DeleteSecretInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.Secret, _ *svcsdk.DeleteSecretOutput, err error) error {
	return err
}
func nopPreUpdate(context.Context, *svcapitypes.Secret, *svcsdk.UpdateSecretInput) error {
	return nil
}
func nopPostUpdate(_ context.Context, _ *svcapitypes.Secret, _ *svcsdk.UpdateSecretOutput, upd managed.ExternalUpdate, err error) (managed.ExternalUpdate, error) {
	return upd, err
}
