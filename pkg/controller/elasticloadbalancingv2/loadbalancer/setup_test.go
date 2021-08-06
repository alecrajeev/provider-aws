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
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
)

var (
	loadBalancerName = "ccloadbalancer"

	testLoadBalancer = svcsdk.LoadBalancer{
		LoadBalancerName: &loadBalancerName,
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
		resultBool   bool
		resultString string
		err          error
	}

	cases := map[string]struct {
		args
		want
	}{
		"NilSourceNoUpdate": {
			args: args{
				cr: loadBalancer(withSpec(v1alpha1.LoadBalancerParameters{})),
				obj: &svcsdk.DescribeLoadBalancersOutput{LoadBalancers: []*svcsdk.LoadBalancer{
					&testLoadBalancer,
				}},
			},
			want: want{
				resultBool:   true,
				resultString: "",
				err:          nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resultBool, resultString := isUpToDateSecurityGroups(tc.args.cr, tc.args.obj)

			// Assert
			if diff := cmp.Diff(tc.want.resultString, resultString, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}

			if resultBool != tc.want.resultBool {
				t.Errorf("r: -want, +got:\n%s", strconv.FormatBool(resultBool))
			}
		})
	}
}
