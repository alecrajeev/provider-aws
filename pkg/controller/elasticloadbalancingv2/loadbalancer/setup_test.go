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
			{SubnetId: aws.String("subnet-000000")},
			{SubnetId: aws.String("subnet-111111")},
		},
	}

	testLoadBalancerDoubleSubnetMappings = svcsdk.LoadBalancer{
		AvailabilityZones: []*svcsdk.AvailabilityZone{
			{SubnetId: aws.String("subnet-000000"),
				LoadBalancerAddresses: []*svcsdk.LoadBalancerAddress{{PrivateIPv4Address: aws.String("172.16.0.6")}}},
			{SubnetId: aws.String("subnet-111111"),
				LoadBalancerAddresses: []*svcsdk.LoadBalancerAddress{{PrivateIPv4Address: aws.String("172.16.24.6")}}},
		},
	}

	testEmptyTags = []*svcsdk.TagDescription{
		{ResourceArn: aws.String("arn:aws:elasticloadbalancing:us-east-2:123456789012:loadbalancer/app/my-load-balancer1/1234567890123456"), Tags: []*svcsdk.Tag{}},
	}

	testExistingTag = []*svcsdk.TagDescription{
		{ResourceArn: aws.String("arn:aws:elasticloadbalancing:us-east-2:123456789012:loadbalancer/app/my-load-balancer2/1234567890123456"), Tags: []*svcsdk.Tag{
			{Key: aws.String("k2"), Value: aws.String("exists_in_obj")},
		}},
	}

	testAddTagsMap = map[string]*string{
		"k1": aws.String("val1"),
	}

	testDescribeTagsOutput = svcsdk.DescribeTagsOutput{
		TagDescriptions: testExistingTag,
	}

	testCRTags = []*v1alpha1.Tag{
		{Key: aws.String("k1"), Value: aws.String("val1")},
	}
)

type args struct {
	cr  *v1alpha1.LoadBalancer
	obj *svcsdk.DescribeLoadBalancersOutput
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
		"NilSourceWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnetMappings}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NilAwsWithUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{
						{SubnetID: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")},
						{SubnetID: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.20.6")}}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerEmptySubnets}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NeedsUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{
						{SubnetID: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")},
						{SubnetID: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.24.6")},
						{SubnetID: aws.String("subnet-222222"), PrivateIPv4Address: aws.String("172.16.28.6")}}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnetMappings}},
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
		"NoUpdateNeededSortOrderIsDifferent": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{
						{SubnetID: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.24.6")},
						{SubnetID: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")}}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnetMappings}},
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"NoUpdateNeeded": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{
						{SubnetID: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")},
						{SubnetID: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.24.6")}}})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancerDoubleSubnetMappings}},
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

func TestDiffTags(t *testing.T) {
	type args struct {
		cr  []*v1alpha1.Tag
		obj *svcsdk.DescribeTagsOutput
	}

	type want struct {
		addTags    map[string]*string
		removeTags []*string
	}

	cases := map[string]struct {
		args
		want
	}{
		"EmptyTags": {
			args: args{
				cr: []*v1alpha1.Tag{},
				obj: &svcsdk.DescribeTagsOutput{
					TagDescriptions: testEmptyTags,
				},
			},
			want: want{
				addTags:    map[string]*string{},
				removeTags: []*string{},
			},
		},
		"AddNewTag": {
			args: args{
				cr: []*v1alpha1.Tag{
					{Key: aws.String("k1"), Value: aws.String("exists_in_cr")},
					{Key: aws.String("k2"), Value: aws.String("exists_in_obj")},
				},
				obj: &svcsdk.DescribeTagsOutput{
					TagDescriptions: testExistingTag},
			},
			want: want{
				addTags: map[string]*string{
					"k1": aws.String("exists_in_cr"),
				},
				removeTags: []*string{},
			},
		},
		"RemoveExistingTag": {
			args: args{
				cr: []*v1alpha1.Tag{},
				obj: &svcsdk.DescribeTagsOutput{
					TagDescriptions: testExistingTag},
			},
			want: want{
				addTags:    map[string]*string{},
				removeTags: []*string{aws.String("k2")},
			},
		},
		"AddAndRemoveWhenKeyChanges": {
			args: args{
				cr: []*v1alpha1.Tag{
					{Key: aws.String("k2"), Value: aws.String("same_key_different_value1")},
				},
				obj: &svcsdk.DescribeTagsOutput{
					TagDescriptions: testExistingTag},
			},
			want: want{
				addTags: map[string]*string{
					"k2": aws.String("same_key_different_value1"),
				},
				removeTags: []*string{
					aws.String("k2")},
			},
		},
		"NoChange": {
			args: args{
				cr: []*v1alpha1.Tag{
					{Key: aws.String("k2"), Value: aws.String("exists_in_obj")},
				},
				obj: &svcsdk.DescribeTagsOutput{
					TagDescriptions: testExistingTag,
				},
			},
			want: want{
				addTags:    map[string]*string{},
				removeTags: []*string{},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			addTags, removeTags, _ := diffTags(tc.args.cr, tc.args.obj)

			// Assert
			if diff := cmp.Diff(tc.want.addTags, addTags, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.removeTags, removeTags, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateSetIPAddressTypeInput(t *testing.T) {
	type args struct {
		cr *v1alpha1.LoadBalancer
	}
	type want struct {
		obj *svcsdk.SetIpAddressTypeInput
	}

	cases := map[string]struct {
		args
		want
	}{
		"ChangeToDualStack": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					IPAddressType: aws.String("dualstack"),
				})),
			},
			want: want{
				obj: &svcsdk.SetIpAddressTypeInput{
					LoadBalancerArn: aws.String(""),
					IpAddressType:   aws.String("dualstack"),
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateSetIPAddressTypeInput(tc.args.cr)

			// Assert
			if diff := cmp.Diff(tc.want.obj, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateSetSecurityGroupsInput(t *testing.T) {
	type args struct {
		cr *v1alpha1.LoadBalancer
	}
	type want struct {
		obj *svcsdk.SetSecurityGroupsInput
	}

	cases := map[string]struct {
		args
		want
	}{
		"UpdateSecurityGroups": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SecurityGroups: []*string{aws.String("sg-111111")},
				})),
			},
			want: want{
				obj: &svcsdk.SetSecurityGroupsInput{
					LoadBalancerArn: aws.String(""),
					SecurityGroups:  []*string{aws.String("sg-111111")},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateSetSecurityGroupsInput(tc.args.cr)

			// Assert
			if diff := cmp.Diff(tc.want.obj, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateSetSubnetsInput(t *testing.T) {
	type args struct {
		cr *v1alpha1.LoadBalancer
	}
	type want struct {
		obj *svcsdk.SetSubnetsInput
	}

	cases := map[string]struct {
		args
		want
	}{
		"UpdateSubnets": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Subnets: []*string{aws.String("subnet-000000"), aws.String("subnet-111111")},
				})),
			},
			want: want{
				obj: &svcsdk.SetSubnetsInput{
					LoadBalancerArn: aws.String(""),
					Subnets:         []*string{aws.String("subnet-000000"), aws.String("subnet-111111")},
				},
			},
		},
		"UpdateSubnetMappings": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					SubnetMappings: []*v1alpha1.SubnetMapping{
						{SubnetID: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")},
						{SubnetID: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.20.6")},
					},
				})),
			},
			want: want{
				obj: &svcsdk.SetSubnetsInput{
					LoadBalancerArn: aws.String(""),
					SubnetMappings: []*svcsdk.SubnetMapping{
						{SubnetId: aws.String("subnet-000000"), PrivateIPv4Address: aws.String("172.16.0.6")},
						{SubnetId: aws.String("subnet-111111"), PrivateIPv4Address: aws.String("172.16.20.6")},
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateSetSubnetsInput(tc.args.cr)

			// Assert
			if diff := cmp.Diff(tc.want.obj, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateAddTagsInput(t *testing.T) {
	type args struct {
		addTags map[string]*string
		cr      *v1alpha1.LoadBalancer
	}
	type want struct {
		obj *svcsdk.AddTagsInput
	}

	cases := map[string]struct {
		args
		want
	}{
		"AddNewTag": {
			args: args{
				addTags: testAddTagsMap,
				cr:      loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
			},
			want: want{
				obj: &svcsdk.AddTagsInput{
					ResourceArns: []*string{aws.String("")},
					Tags: []*svcsdk.Tag{
						{Key: aws.String("k1"), Value: aws.String("val1")},
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateAddTagsInput(tc.args.addTags, tc.args.cr)

			// Assert
			if diff := cmp.Diff(tc.want.obj, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateRemoveTagsInput(t *testing.T) {
	type args struct {
		removeTags []*string
		cr         *v1alpha1.LoadBalancer
	}
	type want struct {
		obj *svcsdk.RemoveTagsInput
	}

	cases := map[string]struct {
		args
		want
	}{
		"AddNewTag": {
			args: args{
				removeTags: []*string{aws.String("k1")},
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{
					Tags: []*v1alpha1.Tag{
						{Key: aws.String("k1"), Value: aws.String("v1")}}})),
			},
			want: want{
				obj: &svcsdk.RemoveTagsInput{
					ResourceArns: []*string{aws.String("")},
					TagKeys:      []*string{aws.String("k1")},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateRemoveTagsInput(tc.args.removeTags, tc.args.cr)

			// Assert
			if diff := cmp.Diff(tc.want.obj, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateMapFromTagsResponseOutput(t *testing.T) {
	type args struct {
		respTags *svcsdk.DescribeTagsOutput
	}
	type want struct {
		objTags map[string]*string
	}

	cases := map[string]struct {
		args
		want
	}{
		"GetMapFromTag": {
			args: args{
				respTags: &testDescribeTagsOutput,
			},
			want: want{
				objTags: map[string]*string{
					"k2": aws.String("exists_in_obj"),
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateMapFromTagsResponseOutput(tc.args.respTags)

			// Assert
			if diff := cmp.Diff(tc.want.objTags, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateMapFromTagsCR(t *testing.T) {
	type args struct {
		crTags []*v1alpha1.Tag
	}
	type want struct {
		objTags map[string]*string
	}

	cases := map[string]struct {
		args
		want
	}{
		"GetMapFromTag": {
			args: args{
				crTags: testCRTags,
			},
			want: want{
				objTags: map[string]*string{
					"k1": aws.String("val1"),
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Act
			actual := GenerateMapFromTagsCR(tc.args.crTags)

			// Assert
			if diff := cmp.Diff(tc.want.objTags, actual, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
