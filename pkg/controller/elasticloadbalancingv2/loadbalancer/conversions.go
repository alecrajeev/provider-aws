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
	svcsdk "github.com/aws/aws-sdk-go/service/elbv2"

	meta "github.com/crossplane/crossplane-runtime/pkg/meta"

	svcapitypes "github.com/crossplane/provider-aws/apis/elasticloadbalancingv2/v1alpha1"
)

// GenerateDescribeTagsInput returns input for read tags
// operation.
func GenerateDescribeTagsInput(cr *svcapitypes.LoadBalancer) *svcsdk.DescribeTagsInput {
	res := &svcsdk.DescribeTagsInput{}
	if len(*cr.Spec.ForProvider.Name) > 0 {
		resourceArn := meta.GetExternalName(cr)
		res.SetResourceArns([]*string{&resourceArn})
	}

	return res
}

// GenerateLoadBalancerWithTags returns the current state in the form of *svcapitypes.LoadBalancer.
// Includes tags information from the DescribeTagsOutput API response.
// Edited from generated GenerateLoadBalancerWithTags to include tags.
// nolint:gocyclo,gosimple,staticcheck
func GenerateLoadBalancerWithTags(resp *svcsdk.DescribeLoadBalancersOutput, respTags *svcsdk.DescribeTagsOutput) *svcapitypes.LoadBalancer {
	cr := &svcapitypes.LoadBalancer{}

	found := false
	for _, elem := range resp.LoadBalancers {
		if elem.CustomerOwnedIpv4Pool != nil {
			cr.Spec.ForProvider.CustomerOwnedIPv4Pool = elem.CustomerOwnedIpv4Pool
		} else {
			cr.Spec.ForProvider.CustomerOwnedIPv4Pool = nil
		}
		if elem.IpAddressType != nil {
			cr.Spec.ForProvider.IPAddressType = elem.IpAddressType
		} else {
			cr.Spec.ForProvider.IPAddressType = nil
		}
		if elem.Scheme != nil {
			cr.Spec.ForProvider.Scheme = elem.Scheme
		} else {
			cr.Spec.ForProvider.Scheme = nil
		}
		if elem.SecurityGroups != nil {
			f9 := []*string{}
			for _, f9iter := range elem.SecurityGroups {
				var f9elem string
				f9elem = *f9iter
				f9 = append(f9, &f9elem)
			}
			cr.Spec.ForProvider.SecurityGroups = f9
		} else {
			cr.Spec.ForProvider.SecurityGroups = nil
		}
		if elem.Type != nil {
			cr.Spec.ForProvider.Type = elem.Type
		} else {
			cr.Spec.ForProvider.Type = nil
		}
		found = true
		break
	}
	if !found {
		return cr
	}

	// Sets tags response values from DescribeTagsOutput into
	// the custom resource.
	if len(respTags.TagDescriptions) > 0 {
		tagDescription := respTags.TagDescriptions[0]
		tags := make([]*svcapitypes.Tag, len(tagDescription.Tags))
		for i, t := range tagDescription.Tags {
			tags[i] = &svcapitypes.Tag{
				Key:   t.Key,
				Value: t.Value,
			}
		}
		cr.Spec.ForProvider.Tags = tags
	}

	return cr
}

// GenerateMapFromTagsResponseOutput returns a map with an input as a tag key, and a value
// as a tag value. The function's input is the response from the DescribeTagsOutput API call.
func GenerateMapFromTagsResponseOutput(respTags *svcsdk.DescribeTagsOutput) map[string]*string {
	var addMap map[string]*string
	tagDescription := respTags.TagDescriptions

	if len(tagDescription) > 0 {
		addMap = make(map[string]*string, len(tagDescription[0].Tags))
		for _, t := range tagDescription[0].Tags {
			addMap[*t.Key] = t.Value
		}
	}

	return addMap
}

// GenerateMapFromTagsCR returns a map with an input as a tag key, and a value
// as a tag value. The function's input is the Tags spec from the LoadBalancer custom resource.
func GenerateMapFromTagsCR(tagsSpec []*svcapitypes.Tag) map[string]*string {
	addMap := make(map[string]*string, len(tagsSpec))
	for _, t := range tagsSpec {
		addMap[*t.Key] = t.Value
	}

	return addMap
}

// GenerateRemoveTagsInput returns the input required to remove tags
// on a LoadBalancer.
func GenerateRemoveTagsInput(removeTags []*string, cr *svcapitypes.LoadBalancer) *svcsdk.RemoveTagsInput {
	loadBalancerArn := meta.GetExternalName(cr)

	return &svcsdk.RemoveTagsInput{
		ResourceArns: []*string{&loadBalancerArn},
		TagKeys:      removeTags,
	}
}

// GenerateAddTagsInput returns the input required to add tags
// on a LoadBalancer
func GenerateAddTagsInput(addTags map[string]*string, cr *svcapitypes.LoadBalancer) *svcsdk.AddTagsInput {
	loadBalanceArn := meta.GetExternalName(cr)
	tags := make([]*svcsdk.Tag, len(addTags))
	i := 0
	for k, v := range addTags {
		key := k
		tag := svcsdk.Tag{
			Key:   &key,
			Value: v,
		}
		tags[i] = &tag
		i++
	}

	return &svcsdk.AddTagsInput{
		ResourceArns: []*string{&loadBalanceArn},
		Tags:         tags,
	}
}
