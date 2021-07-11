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

package acm

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	awsacm "github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"
	"github.com/alecrajeev/crossplane-runtime/pkg/resource"
	"github.com/alecrajeev/crossplane-runtime/pkg/test"

	v1alpha1 "github.com/crossplane/provider-aws/apis/acm/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	acm "github.com/crossplane/provider-aws/pkg/clients/acm"
	"github.com/crossplane/provider-aws/pkg/clients/acm/fake"
)

var (
	// an arbitrary managed resource
	unexpecedItem  resource.Managed
	domainName     = "some.site"
	certificateArn = "somearn"

	errBoom = errors.New("boom")
)

type args struct {
	acm acm.Client
	cr  resource.Managed
}

type certificateModifier func(*v1alpha1.Certificate)

func withConditions(c ...xpv1.Condition) certificateModifier {
	return func(r *v1alpha1.Certificate) { r.Status.ConditionedStatus.Conditions = c }
}

func withDomainName() certificateModifier {
	return func(r *v1alpha1.Certificate) {
		r.Spec.ForProvider.DomainName = domainName
		meta.SetExternalName(r, certificateArn)
	}
}

func withCertificateTransparencyLoggingPreference() certificateModifier {
	certificateTransparencyLoggingPreference := awsacm.CertificateTransparencyLoggingPreferenceDisabled

	return func(r *v1alpha1.Certificate) {
		r.Spec.ForProvider.CertificateTransparencyLoggingPreference = &certificateTransparencyLoggingPreference
		meta.SetExternalName(r, certificateArn)
	}
}

func withTags() certificateModifier {
	return func(r *v1alpha1.Certificate) {
		r.Spec.ForProvider.Tags = append(r.Spec.ForProvider.Tags, v1alpha1.Tag{
			Key:   "Name",
			Value: "somename",
		})
		meta.SetExternalName(r, certificateArn)
	}
}

func withCertificateArn() certificateModifier {
	return func(r *v1alpha1.Certificate) {
		certificateTransparencyLoggingPreference := awsacm.CertificateTransparencyLoggingPreferenceDisabled

		r.Status.AtProvider.CertificateARN = certificateArn
		r.Spec.ForProvider.CertificateTransparencyLoggingPreference = &certificateTransparencyLoggingPreference
		meta.SetExternalName(r, certificateArn)
	}
}

func withStatus(status awsacm.CertificateStatus) certificateModifier {
	return func(r *v1alpha1.Certificate) {
		r.Status.AtProvider.Status = status
	}
}

func certificate(m ...certificateModifier) *v1alpha1.Certificate {
	cr := &v1alpha1.Certificate{}
	meta.SetExternalName(cr, certificateArn)
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestObserve(t *testing.T) {

	type want struct {
		cr     resource.Managed
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"ValidInput": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDescribeCertificateRequest: func(input *awsacm.DescribeCertificateInput) awsacm.DescribeCertificateRequest {
						return awsacm.DescribeCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.DescribeCertificateOutput{
								Certificate: &awsacm.CertificateDetail{
									CertificateArn: aws.String(certificateArn),
									Options:        &awsacm.CertificateOptions{CertificateTransparencyLoggingPreference: awsacm.CertificateTransparencyLoggingPreferenceDisabled},
									Status:         awsacm.CertificateStatusIssued,
								},
							}},
						}
					},
					MockListTagsForCertificateRequest: func(input *awsacm.ListTagsForCertificateInput) awsacm.ListTagsForCertificateRequest {
						return awsacm.ListTagsForCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.ListTagsForCertificateOutput{
								Tags: []awsacm.Tag{{}},
							}},
						}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr: certificate(withCertificateArn(), withStatus(awsacm.CertificateStatusIssued), withConditions(xpv1.Available())),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
				},
			},
		},
		"InValidInput": {
			args: args{
				cr: unexpecedItem,
			},
			want: want{
				cr:  unexpecedItem,
				err: errors.New(errUnexpectedObject),
			},
		},
		"ClientError": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDescribeCertificateRequest: func(input *awsacm.DescribeCertificateInput) awsacm.DescribeCertificateRequest {
						return awsacm.DescribeCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom},
						}
					},
				},
				cr: certificate(withCertificateArn()),
			},
			want: want{
				cr:  certificate(withCertificateArn()),
				err: awsclient.Wrap(errBoom, errGet),
			},
		},
		"ResourceDoesNotExist": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDescribeCertificateRequest: func(input *awsacm.DescribeCertificateInput) awsacm.DescribeCertificateRequest {
						return awsacm.DescribeCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: awserr.New(awsacm.ErrCodeResourceNotFoundException, "", nil)},
						}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr: certificate(),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{
				client: tc.acm,
				kube: &test.MockClient{
					MockUpdate: test.NewMockUpdateFn(nil),
				},
			}
			o, err := e.Observe(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}

			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {

	type want struct {
		cr     resource.Managed
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"ValidInput": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockRequestCertificateRequest: func(input *awsacm.RequestCertificateInput) awsacm.RequestCertificateRequest {
						return awsacm.RequestCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RequestCertificateOutput{
								CertificateArn: aws.String(certificateArn),
							}},
						}
					},
				},
				cr: certificate(withDomainName()),
			},
			want: want{
				cr: certificate(
					withDomainName()),
				result: managed.ExternalCreation{ExternalNameAssigned: true},
			},
		},
		"InValidInput": {
			args: args{
				cr: unexpecedItem,
			},
			want: want{
				cr:  unexpecedItem,
				err: errors.New(errUnexpectedObject),
			},
		},
		"ClientError": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockRequestCertificateRequest: func(input *awsacm.RequestCertificateInput) awsacm.RequestCertificateRequest {
						return awsacm.RequestCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom},
						}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr:  certificate(),
				err: awsclient.Wrap(errBoom, errCreate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{
				client: tc.acm,
				kube: &test.MockClient{
					MockUpdate: test.NewMockUpdateFn(nil),
				},
			}
			o, err := e.Create(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {

	type want struct {
		cr     resource.Managed
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"ValidInput": {
			args: args{
				acm: &fake.MockCertificateClient{

					MockUpdateCertificateOptionsRequest: func(input *awsacm.UpdateCertificateOptionsInput) awsacm.UpdateCertificateOptionsRequest {
						return awsacm.UpdateCertificateOptionsRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.UpdateCertificateOptionsOutput{}}}
					},
					MockListTagsForCertificateRequest: func(input *awsacm.ListTagsForCertificateInput) awsacm.ListTagsForCertificateRequest {
						return awsacm.ListTagsForCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.ListTagsForCertificateOutput{
								Tags: []awsacm.Tag{{}},
							}},
						}
					},
					MockRemoveTagsFromCertificateRequest: func(input *awsacm.RemoveTagsFromCertificateInput) awsacm.RemoveTagsFromCertificateRequest {
						return awsacm.RemoveTagsFromCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RemoveTagsFromCertificateOutput{}}}
					},
					MockAddTagsToCertificateRequest: func(input *awsacm.AddTagsToCertificateInput) awsacm.AddTagsToCertificateRequest {
						return awsacm.AddTagsToCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.AddTagsToCertificateOutput{}}}
					},
					MockRenewCertificateRequest: func(input *awsacm.RenewCertificateInput) awsacm.RenewCertificateRequest {
						return awsacm.RenewCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RenewCertificateOutput{}}}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr: certificate(),
			},
		},
		"InValidInput": {
			args: args{
				cr: unexpecedItem,
			},
			want: want{
				cr:  unexpecedItem,
				err: errors.New(errUnexpectedObject),
			},
		},
		"ClientUpdateCertificateOptionsError": {
			args: args{
				acm: &fake.MockCertificateClient{

					MockUpdateCertificateOptionsRequest: func(input *awsacm.UpdateCertificateOptionsInput) awsacm.UpdateCertificateOptionsRequest {
						return awsacm.UpdateCertificateOptionsRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom}}
					},
					MockListTagsForCertificateRequest: func(input *awsacm.ListTagsForCertificateInput) awsacm.ListTagsForCertificateRequest {
						return awsacm.ListTagsForCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.ListTagsForCertificateOutput{
								Tags: []awsacm.Tag{{}},
							}},
						}
					},
					MockRemoveTagsFromCertificateRequest: func(input *awsacm.RemoveTagsFromCertificateInput) awsacm.RemoveTagsFromCertificateRequest {
						return awsacm.RemoveTagsFromCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RemoveTagsFromCertificateOutput{}}}
					},
					MockAddTagsToCertificateRequest: func(input *awsacm.AddTagsToCertificateInput) awsacm.AddTagsToCertificateRequest {
						return awsacm.AddTagsToCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.AddTagsToCertificateOutput{}}}
					},
					MockRenewCertificateRequest: func(input *awsacm.RenewCertificateInput) awsacm.RenewCertificateRequest {
						return awsacm.RenewCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RenewCertificateOutput{}}}
					},
				},
				cr: certificate(withCertificateTransparencyLoggingPreference()),
			},
			want: want{
				cr:  certificate(withCertificateTransparencyLoggingPreference()),
				err: awsclient.Wrap(errBoom, errUpdate),
			},
		},
		"ClientUpdateTagsError": {
			args: args{
				acm: &fake.MockCertificateClient{

					MockUpdateCertificateOptionsRequest: func(input *awsacm.UpdateCertificateOptionsInput) awsacm.UpdateCertificateOptionsRequest {
						return awsacm.UpdateCertificateOptionsRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.UpdateCertificateOptionsOutput{}}}
					},
					MockListTagsForCertificateRequest: func(input *awsacm.ListTagsForCertificateInput) awsacm.ListTagsForCertificateRequest {
						return awsacm.ListTagsForCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.ListTagsForCertificateOutput{
								Tags: []awsacm.Tag{{}},
							}},
						}
					},
					MockRemoveTagsFromCertificateRequest: func(input *awsacm.RemoveTagsFromCertificateInput) awsacm.RemoveTagsFromCertificateRequest {
						return awsacm.RemoveTagsFromCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RemoveTagsFromCertificateOutput{}}}
					},
					MockAddTagsToCertificateRequest: func(input *awsacm.AddTagsToCertificateInput) awsacm.AddTagsToCertificateRequest {
						return awsacm.AddTagsToCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom}}
					},
					MockRenewCertificateRequest: func(input *awsacm.RenewCertificateInput) awsacm.RenewCertificateRequest {
						return awsacm.RenewCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.RenewCertificateOutput{}}}
					},
				},
				cr: certificate(withTags()),
			},
			want: want{
				cr:  certificate(withTags()),
				err: awsclient.Wrap(errBoom, errAddTagsFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.acm}
			o, err := e.Update(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	type want struct {
		cr  resource.Managed
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"ValidInput": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDeleteCertificateRequest: func(input *awsacm.DeleteCertificateInput) awsacm.DeleteCertificateRequest {
						return awsacm.DeleteCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsacm.DeleteCertificateOutput{}},
						}
					},
				},
				cr: certificate(withCertificateTransparencyLoggingPreference()),
			},
			want: want{
				cr: certificate(withCertificateTransparencyLoggingPreference(),
					withConditions(xpv1.Deleting())),
			},
		},
		"InValidInput": {
			args: args{
				cr: unexpecedItem,
			},
			want: want{
				cr:  unexpecedItem,
				err: errors.New(errUnexpectedObject),
			},
		},
		"ClientError": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDeleteCertificateRequest: func(input *awsacm.DeleteCertificateInput) awsacm.DeleteCertificateRequest {
						return awsacm.DeleteCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom},
						}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr:  certificate(withConditions(xpv1.Deleting())),
				err: awsclient.Wrap(errBoom, errDelete),
			},
		},
		"ResourceDoesNotExist": {
			args: args{
				acm: &fake.MockCertificateClient{
					MockDeleteCertificateRequest: func(input *awsacm.DeleteCertificateInput) awsacm.DeleteCertificateRequest {
						return awsacm.DeleteCertificateRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: awserr.New(awsacm.ErrCodeResourceNotFoundException, "", nil)},
						}
					},
				},
				cr: certificate(),
			},
			want: want{
				cr: certificate(withConditions(xpv1.Deleting())),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.acm}
			err := e.Delete(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
