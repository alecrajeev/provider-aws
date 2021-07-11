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

package iamuser

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
	userName      = "some user"

	errBoom = errors.New("boom")
)

type args struct {
	iam iam.UserClient
	cr  resource.Managed
}

type userModifier func(*v1alpha1.IAMUser)

func withConditions(c ...xpv1.Condition) userModifier {
	return func(r *v1alpha1.IAMUser) { r.Status.ConditionedStatus.Conditions = c }
}

func withExternalName(name string) userModifier {
	return func(r *v1alpha1.IAMUser) { meta.SetExternalName(r, name) }
}

func user(m ...userModifier) *v1alpha1.IAMUser {
	cr := &v1alpha1.IAMUser{}
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
		"VaildInput": {
			args: args{
				iam: &fake.MockUserClient{
					MockGetUser: func(input *awsiam.GetUserInput) awsiam.GetUserRequest {
						return awsiam.GetUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.GetUserOutput{
								User: &awsiam.User{},
							}},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr: user(withExternalName(userName),
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
		"GetUserError": {
			args: args{
				iam: &fake.MockUserClient{
					MockGetUser: func(input *awsiam.GetUserInput) awsiam.GetUserRequest {
						return awsiam.GetUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr:  user(withExternalName(userName)),
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
		"VaildInput": {
			args: args{
				iam: &fake.MockUserClient{
					MockCreateUser: func(input *awsiam.CreateUserInput) awsiam.CreateUserRequest {
						return awsiam.CreateUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.CreateUserOutput{}},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr: user(
					withExternalName(userName),
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
				iam: &fake.MockUserClient{
					MockCreateUser: func(input *awsiam.CreateUserInput) awsiam.CreateUserRequest {
						return awsiam.CreateUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: user(),
			},
			want: want{
				cr:  user(withConditions(xpv1.Creating())),
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
		"VaildInput": {
			args: args{
				iam: &fake.MockUserClient{
					MockUpdateUser: func(input *awsiam.UpdateUserInput) awsiam.UpdateUserRequest {
						return awsiam.UpdateUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.UpdateUserOutput{}},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr: user(withExternalName(userName)),
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
		"VaildInput": {
			args: args{
				iam: &fake.MockUserClient{
					MockDeleteUser: func(input *awsiam.DeleteUserInput) awsiam.DeleteUserRequest {
						return awsiam.DeleteUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsiam.DeleteUserOutput{}},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr: user(withExternalName(userName),
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
				iam: &fake.MockUserClient{
					MockDeleteUser: func(input *awsiam.DeleteUserInput) awsiam.DeleteUserRequest {
						return awsiam.DeleteUserRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: user(withExternalName(userName)),
			},
			want: want{
				cr: user(withExternalName(userName),
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
