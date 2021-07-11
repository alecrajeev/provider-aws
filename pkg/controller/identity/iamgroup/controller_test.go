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

package iamgroup

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/alecrajeev/crossplane-runtime/pkg/meta"
	"github.com/alecrajeev/crossplane-runtime/pkg/reconciler/managed"
	"github.com/alecrajeev/crossplane-runtime/pkg/resource"

	xpv1 "github.com/alecrajeev/crossplane-runtime/apis/common/v1"
	"github.com/alecrajeev/crossplane-runtime/pkg/test"

	"github.com/crossplane/provider-aws/apis/identity/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/iam"
	"github.com/crossplane/provider-aws/pkg/clients/iam/fake"
)

var (
	unexpecedItem resource.Managed
	groupName     = "some group"

	errBoom = errors.New("boom")
)

const (
	groupPath = "group-path"
)

type args struct {
	iam iam.GroupClient
	cr  resource.Managed
}

type groupModifier func(*v1alpha1.IAMGroup)

func withConditions(c ...xpv1.Condition) groupModifier {
	return func(r *v1alpha1.IAMGroup) { r.Status.ConditionedStatus.Conditions = c }
}

func withExternalName(name string) groupModifier {
	return func(r *v1alpha1.IAMGroup) { meta.SetExternalName(r, name) }
}

func withGroupPath(groupPath string) groupModifier {
	return func(r *v1alpha1.IAMGroup) { r.Spec.ForProvider.Path = &groupPath }
}

func group(m ...groupModifier) *v1alpha1.IAMGroup {
	cr := &v1alpha1.IAMGroup{}
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
				iam: &fake.MockGroupClient{
					MockGetGroup: func(input *awsiam.GetGroupInput) awsiam.GetGroupRequest {
						return awsiam.GetGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.GetGroupOutput{
								Group: &awsiam.Group{
									GroupName: aws.String(groupName),
									Path:      aws.String(groupPath),
								},
							}},
						}
					},
				},
				cr: group(withExternalName(groupName),
					withGroupPath(groupPath)),
			},
			want: want{
				cr: group(withExternalName(groupName),
					withGroupPath(groupPath),
					withConditions(xpv1.Available())),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
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
		"GetGroupError": {
			args: args{
				iam: &fake.MockGroupClient{
					MockGetGroup: func(input *awsiam.GetGroupInput) awsiam.GetGroupRequest {
						return awsiam.GetGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom, Retryer: aws.NoOpRetryer{}},
						}
					},
				},
				cr: group(withExternalName(groupName)),
			},
			want: want{
				cr:  group(withExternalName(groupName)),
				err: awsclient.Wrap(errBoom, errGet),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.iam}
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
				iam: &fake.MockGroupClient{
					MockCreateGroup: func(input *awsiam.CreateGroupInput) awsiam.CreateGroupRequest {
						return awsiam.CreateGroupRequest{Request: &aws.Request{HTTPRequest: &http.Request{}, Data: &awsiam.CreateGroupOutput{}, Retryer: aws.NoOpRetryer{}}}
					},
				},
				cr: group(withExternalName(groupName)),
			},
			want: want{
				cr: group(
					withExternalName(groupName),
					withConditions(xpv1.Creating())),
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
				iam: &fake.MockGroupClient{
					MockCreateGroup: func(input *awsiam.CreateGroupInput) awsiam.CreateGroupRequest {
						return awsiam.CreateGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom},
						}
					},
				},
				cr: group(),
			},
			want: want{
				cr:  group(withConditions(xpv1.Creating())),
				err: awsclient.Wrap(errBoom, errCreate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.iam}
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
				iam: &fake.MockGroupClient{
					MockUpdateGroup: func(input *awsiam.UpdateGroupInput) awsiam.UpdateGroupRequest {
						return awsiam.UpdateGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.UpdateGroupOutput{}},
						}
					},
				},
				cr: group(withExternalName(groupName)),
			},
			want: want{
				cr: group(withExternalName(groupName)),
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
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.iam}
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
				iam: &fake.MockGroupClient{
					MockDeleteGroup: func(input *awsiam.DeleteGroupInput) awsiam.DeleteGroupRequest {
						return awsiam.DeleteGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.DeleteGroupOutput{}},
						}
					},
				},
				cr: group(withExternalName(groupName)),
			},
			want: want{
				cr: group(withExternalName(groupName),
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
		"DeleteError": {
			args: args{
				iam: &fake.MockGroupClient{
					MockDeleteGroup: func(input *awsiam.DeleteGroupInput) awsiam.DeleteGroupRequest {
						return awsiam.DeleteGroupRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: errBoom},
						}
					},
				},
				cr: group(withExternalName(groupName)),
			},
			want: want{
				cr: group(withExternalName(groupName),
					withConditions(xpv1.Deleting())),
				err: awsclient.Wrap(errBoom, errDelete),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.iam}
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
