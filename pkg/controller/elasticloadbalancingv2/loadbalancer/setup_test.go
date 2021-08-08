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

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
)

var (
	testLoadBalancerNilSecurityGroups = svcsdk.LoadBalancer{
		LoadBalancerName: aws.String("testloadbalancer"),
	}

	testLoadBalancerEmptySecurityGroups = svcsdk.LoadBalancer{
		SecurityGroups: []*string{},
	}

	testLoadBalancerSingleSecurityGroup = svcsdk.LoadBalancer{
		SecurityGroups: []*string{aws.String("sg-000000")},
	}

	testLoadBalancerDoubleSecurityGroups = svcsdk.LoadBalancer{
		SecurityGroups: []*string{aws.String("sg-000000"), aws.String("sg-111111")},
	}

	testLoadBalancerEmptySubnets = svcsdk.LoadBalancer{
		AvailabilityZones: []*svcsdk.AvailabilityZone{},
	}

	testLoadBalancerDoubleSubnets = svcsdk.LoadBalancer{
		AvailabilityZones: []*svcsdk.AvailabilityZone{
			&svcsdk.AvailabilityZone{SubnetId: aws.String("subnet-000000")},
			&svcsdk.AvailabilityZone{SubnetId: aws.String("subnet-111111")},
		},
	}
)

type args struct {
	cr      *v1alpha1.LoadBalancer
	obj     *svcsdk.DescribeLoadBalancersOutput
	objTags *svcsdk.DescribeTagsOutput
}

type loadBalancerModifier func(*v1alpha1.LoadBalancer)

func withSpec(p v1alpha1.LoadBalancerParameters) loadBalancerModifier {
	return func(r *v1alpha1.LoadBalancer) { r.Spec.ForProvider = p }
}

func loadBalancer(m ...loadBalancerModifier) *v1alpha1.LoadBalancer {
	cr := &v1alpha1.LoadBalancer{}
	cr.Name = "test-load-balancer"
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestIsUpToDateSecurityGroups(t *testing.T) {
	type want struct {
		result bool
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"NilSourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerNilSecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NilSourceNilAwsNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"EmptySourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NilSourceWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerSingleSecurityGroup}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NilAwsWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{aws.String("sg-000000")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySecurityGroups,
				}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NeedsUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{aws.String("sg-111111")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerSingleSecurityGroup}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NoUpdateNeededSortOrderIsDifferent": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{aws.String("sg-111111"), aws.String("sg-000000")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NoUpdateNeeded": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{aws.String("sg-000000")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerSingleSecurityGroup}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			result, _ := isUpToDateSecurityGroups(tc.args.cr, tc.args.obj)

			// Assert
			if diff := cmp.Diff(tc.want.result, result, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDateSubnets(t *testing.T) {
	type want struct {
		result bool
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"NilSourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerNilSecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NilSourceNilAwsNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"EmptySourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NilSourceWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnets}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NilAwsWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{aws.String("subnet-000000"), aws.String("subnet-111111")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets,
				}},
			},
		},
		"NeedsUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{aws.String("subnet-000000"), aws.String("subnet-111111"), aws.String("subnet-222222")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnets,
				}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NoUpdateNeededSortOrderIsDifferent": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{aws.String("subnet-111111"), aws.String("subnet-000000")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnets,
				}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NoUpdateNeeded": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{aws.String("subnet-000000"), aws.String("subnet-111111")}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnets,
				}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			result, _ := isUpToDateSubnets(tc.args.cr, tc.args.obj)

			// Assert
			if diff := cmp.Diff(tc.want.result, result, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDateSubnetMappings(t *testing.T) {
	type want struct {
		result bool
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"NilSourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerNilSecurityGroups}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NilSourceNilAwsNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"EmptySourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			result, _ := isUpToDateSubnetMappings(tc.args.cr, tc.args.obj)

			// Assert
			if diff := cmp.Diff(tc.want.result, result, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
