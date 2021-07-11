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

package address

import (
	"context"
	"net/http"
	"testing"

	"github.com/alecrajeev/crossplane-runtime/pkg/resource"
	"github.com/alecrajeev/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"

	"github.com/crossplane/provider-aws/apis/ec2/v1beta1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/ec2"
	"github.com/crossplane/provider-aws/pkg/clients/ec2/fake"
)

var (
	allocationID   = "some Id"
	domainVpc      = "vpc"
	domainStandard = "standard"
	publicIP       = "1.1.1.1"
	errBoom        = errors.New("boom")
)

type args struct {
	address ec2.AddressClient
	kube    client.Client
	cr      *v1beta1.Address
}

type addressModifier func(*v1beta1.Address)

func withTags(tagMaps ...map[string]string) addressModifier {
	var tagList []v1beta1.Tag
	for _, tagMap := range tagMaps {
		for k, v := range tagMap {
			tagList = append(tagList, v1beta1.Tag{Key: k, Value: v})
		}
	}
	return func(r *v1beta1.Address) { r.Spec.ForProvider.Tags = tagList }
}

func withExternalName(name string) addressModifier {
	return func(r *v1beta1.Address) { meta.SetExternalName(r, name) }
}

func withConditions(c ...xpv1.Condition) addressModifier {
	return func(r *v1beta1.Address) { r.Status.ConditionedStatus.Conditions = c }
}

func withSpec(p v1beta1.AddressParameters) addressModifier {
	return func(r *v1beta1.Address) { r.Spec.ForProvider = p }
}

func withStatus(s v1beta1.AddressObservation) addressModifier {
	return func(r *v1beta1.Address) { r.Status.AtProvider = s }
}

func address(m ...addressModifier) *v1beta1.Address {
	cr := &v1beta1.Address{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

var _ managed.ExternalClient = &external{}
var _ managed.ExternalConnecter = &connector{}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1beta1.Address
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				kube: &test.MockClient{
					MockUpdate: test.NewMockClient().Update,
				},
				address: &fake.MockAddressClient{
					MockDescribe: func(input *awsec2.DescribeAddressesInput) awsec2.DescribeAddressesRequest {
						return awsec2.DescribeAddressesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.DescribeAddressesOutput{
								Addresses: []awsec2.Address{{
									AllocationId: &allocationID,
								}},
							}},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withExternalName(allocationID)),
			},
			want: want{
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withStatus(v1beta1.AddressObservation{
					AllocationID: allocationID,
				}), withExternalName(allocationID),
					withConditions(xpv1.Available())),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
			},
		},
		"MultipleAddresses": {
			args: args{
				kube: &test.MockClient{
					MockUpdate: test.NewMockClient().Update,
				},
				address: &fake.MockAddressClient{
					MockDescribe: func(input *awsec2.DescribeAddressesInput) awsec2.DescribeAddressesRequest {
						return awsec2.DescribeAddressesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.DescribeAddressesOutput{
								Addresses: []awsec2.Address{{}, {}},
							}},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withExternalName(allocationID)),
			},
			want: want{
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withExternalName(allocationID)),
				err: errors.New(errMultipleItems),
			},
		},
		"DescribeFail": {
			args: args{
				kube: &test.MockClient{
					MockUpdate: test.NewMockClient().Update,
				},
				address: &fake.MockAddressClient{
					MockDescribe: func(input *awsec2.DescribeAddressesInput) awsec2.DescribeAddressesRequest {
						return awsec2.DescribeAddressesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withExternalName(allocationID)),
			},
			want: want{
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				}), withExternalName(allocationID)),
				err: awsclient.Wrap(errBoom, errDescribe),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.address}
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
		cr     *v1beta1.Address
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulVPC": {
			args: args{
				kube: &test.MockClient{
					MockUpdate:       test.NewMockClient().Update,
					MockStatusUpdate: test.NewMockClient().MockStatusUpdate,
				},
				address: &fake.MockAddressClient{
					MockAllocate: func(input *awsec2.AllocateAddressInput) awsec2.AllocateAddressRequest {
						return awsec2.AllocateAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.AllocateAddressOutput{
								AllocationId: &allocationID,
							}},
						}
					},
				},
				cr: address(),
			},
			want: want{
				cr: address(withExternalName(allocationID),
					withConditions(xpv1.Creating())),
				result: managed.ExternalCreation{ExternalNameAssigned: true},
			},
		},
		"SuccessfulStandard": {
			args: args{
				kube: &test.MockClient{
					MockUpdate:       test.NewMockClient().Update,
					MockStatusUpdate: test.NewMockClient().MockStatusUpdate,
				},
				address: &fake.MockAddressClient{
					MockAllocate: func(input *awsec2.AllocateAddressInput) awsec2.AllocateAddressRequest {
						return awsec2.AllocateAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.AllocateAddressOutput{
								PublicIp: &publicIP,
							}},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainStandard,
				})),
			},
			want: want{
				cr: address(withExternalName(publicIP),
					withConditions(xpv1.Creating()),
					withSpec(v1beta1.AddressParameters{
						Domain: &domainStandard,
					})),
				result: managed.ExternalCreation{ExternalNameAssigned: true},
			},
		},
		"CreateFail": {
			args: args{
				kube: &test.MockClient{
					MockUpdate:       test.NewMockClient().Update,
					MockStatusUpdate: test.NewMockClient().MockStatusUpdate,
				},
				address: &fake.MockAddressClient{
					MockAllocate: func(input *awsec2.AllocateAddressInput) awsec2.AllocateAddressRequest {
						return awsec2.AllocateAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: address(),
			},
			want: want{
				cr:  address(withConditions(xpv1.Creating())),
				err: awsclient.Wrap(errBoom, errCreate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.address}
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
		cr     *v1beta1.Address
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				address: &fake.MockAddressClient{

					MockCreateTagsRequest: func(input *awsec2.CreateTagsInput) awsec2.CreateTagsRequest {
						return awsec2.CreateTagsRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.CreateTagsOutput{}},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				})),
			},
			want: want{
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				})),
			},
		},
		"ModifyFailed": {
			args: args{
				address: &fake.MockAddressClient{
					MockCreateTagsRequest: func(input *awsec2.CreateTagsInput) awsec2.CreateTagsRequest {
						return awsec2.CreateTagsRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				})),
			},
			want: want{
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainVpc,
				})),
				err: awsclient.Wrap(errBoom, errCreateTags),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.address}
			u, err := e.Update(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, u); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestRelease(t *testing.T) {
	type want struct {
		cr  *v1beta1.Address
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulVPC": {
			args: args{
				address: &fake.MockAddressClient{
					MockRelease: func(input *awsec2.ReleaseAddressInput) awsec2.ReleaseAddressRequest {
						return awsec2.ReleaseAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.ReleaseAddressOutput{}},
						}
					},
				},
				cr: address(),
			},
			want: want{
				cr: address(withConditions(xpv1.Deleting())),
			},
		},
		"SuccessfulStandard": {
			args: args{
				address: &fake.MockAddressClient{
					MockRelease: func(input *awsec2.ReleaseAddressInput) awsec2.ReleaseAddressRequest {
						return awsec2.ReleaseAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsec2.ReleaseAddressOutput{}},
						}
					},
				},
				cr: address(withSpec(v1beta1.AddressParameters{
					Domain: &domainStandard,
				})),
			},
			want: want{
				cr: address(withConditions(xpv1.Deleting()),
					withSpec(v1beta1.AddressParameters{
						Domain: &domainStandard,
					}),
				),
			},
		},
		"DeleteFailed": {
			args: args{
				address: &fake.MockAddressClient{
					MockRelease: func(input *awsec2.ReleaseAddressInput) awsec2.ReleaseAddressRequest {
						return awsec2.ReleaseAddressRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: address(),
			},
			want: want{
				cr:  address(withConditions(xpv1.Deleting())),
				err: awsclient.Wrap(errBoom, errDelete),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.address}
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

func TestInitialize(t *testing.T) {
	type args struct {
		cr   *v1beta1.Address
		kube client.Client
	}
	type want struct {
		cr  *v1beta1.Address
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				cr:   address(withTags(map[string]string{"foo": "bar"})),
				kube: &test.MockClient{MockUpdate: test.NewMockUpdateFn(nil)},
			},
			want: want{
				cr: address(withTags(resource.GetExternalTags(address()), map[string]string{"foo": "bar"})),
			},
		},
		"UpdateFailed": {
			args: args{
				cr:   address(),
				kube: &test.MockClient{MockUpdate: test.NewMockUpdateFn(errBoom)},
			},
			want: want{
				err: errors.Wrap(errBoom, errKubeUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &tagger{kube: tc.kube}
			err := e.Initialize(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, cmpopts.SortSlices(func(a, b v1beta1.Tag) bool { return a.Key > b.Key })); err == nil && diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
