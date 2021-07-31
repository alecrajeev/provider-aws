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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
)

type Action struct {
	// Request parameters to use when integrating with Amazon Cognito to authenticate
	// users.
	AuthenticateCognitoConfig *AuthenticateCognitoActionConfig `json:"authenticateCognitoConfig,omitempty"`
	// Request parameters when using an identity provider (IdP) that is compliant
	// with OpenID Connect (OIDC) to authenticate users.
	AuthenticateOidcConfig *AuthenticateOidcActionConfig `json:"authenticateOidcConfig,omitempty"`
	// Information about an action that returns a custom HTTP response.
	FixedResponseConfig *FixedResponseActionConfig `json:"fixedResponseConfig,omitempty"`
	// Information about a forward action.
	ForwardConfig *ForwardActionConfig `json:"forwardConfig,omitempty"`

	Order *int64 `json:"order,omitempty"`
	// Information about a redirect action.
	//
	// A URI consists of the following components: protocol://hostname:port/path?query.
	// You must modify at least one of the following components to avoid a redirect
	// loop: protocol, hostname, port, or path. Any components that you do not modify
	// retain their original values.
	//
	// You can reuse URI components using the following reserved keywords:
	//
	//    * #{protocol}
	//
	//    * #{host}
	//
	//    * #{port}
	//
	//    * #{path} (the leading "/" is removed)
	//
	//    * #{query}
	//
	// For example, you can change the path to "/new/#{path}", the hostname to "example.#{host}",
	// or the query to "#{query}&value=xyz".
	RedirectConfig *RedirectActionConfig `json:"redirectConfig,omitempty"`

	TargetGroupARN *string `json:"targetGroupARN,omitempty"`

	Type *string `json:"type_,omitempty"`
}

type AuthenticateCognitoActionConfig struct {
	AuthenticationRequestExtraParams map[string]*string `json:"authenticationRequestExtraParams,omitempty"`

	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest,omitempty"`

	Scope *string `json:"scope,omitempty"`

	SessionCookieName *string `json:"sessionCookieName,omitempty"`

	SessionTimeout *int64 `json:"sessionTimeout,omitempty"`

	UserPoolARN *string `json:"userPoolARN,omitempty"`

	UserPoolClientID *string `json:"userPoolClientID,omitempty"`

	UserPoolDomain *string `json:"userPoolDomain,omitempty"`
}

type AuthenticateOidcActionConfig struct {
	AuthenticationRequestExtraParams map[string]*string `json:"authenticationRequestExtraParams,omitempty"`

	AuthorizationEndpoint *string `json:"authorizationEndpoint,omitempty"`

	ClientID *string `json:"clientID,omitempty"`

	ClientSecret *string `json:"clientSecret,omitempty"`

	Issuer *string `json:"issuer,omitempty"`

	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest,omitempty"`

	Scope *string `json:"scope,omitempty"`

	SessionCookieName *string `json:"sessionCookieName,omitempty"`

	SessionTimeout *int64 `json:"sessionTimeout,omitempty"`

	TokenEndpoint *string `json:"tokenEndpoint,omitempty"`

	UseExistingClientSecret *bool `json:"useExistingClientSecret,omitempty"`

	UserInfoEndpoint *string `json:"userInfoEndpoint,omitempty"`
}

type AvailabilityZone struct {
	LoadBalancerAddresses []*LoadBalancerAddress `json:"loadBalancerAddresses,omitempty"`

	OutpostID *string `json:"outpostID,omitempty"`

	SubnetID *string `json:"subnetID,omitempty"`

	ZoneName *string `json:"zoneName,omitempty"`
}

type Certificate struct {
	CertificateARN *string `json:"certificateARN,omitempty"`

	IsDefault *bool `json:"isDefault,omitempty"`
}

type FixedResponseActionConfig struct {
	ContentType *string `json:"contentType,omitempty"`

	MessageBody *string `json:"messageBody,omitempty"`

	StatusCode *string `json:"statusCode,omitempty"`
}

type ForwardActionConfig struct {
	// Information about the target group stickiness for a rule.
	TargetGroupStickinessConfig *TargetGroupStickinessConfig `json:"targetGroupStickinessConfig,omitempty"`

	TargetGroups []*TargetGroupTuple `json:"targetGroups,omitempty"`
}

type HTTPHeaderConditionConfig struct {
	HTTPHeaderName *string `json:"httpHeaderName,omitempty"`

	Values []*string `json:"values,omitempty"`
}

type HTTPRequestMethodConditionConfig struct {
	Values []*string `json:"values,omitempty"`
}

type HostHeaderConditionConfig struct {
	Values []*string `json:"values,omitempty"`
}

type Listener_SDK struct {
	AlpnPolicy []*string `json:"alpnPolicy,omitempty"`

	Certificates []*Certificate `json:"certificates,omitempty"`

	DefaultActions []*Action `json:"defaultActions,omitempty"`

	ListenerARN *string `json:"listenerARN,omitempty"`

	LoadBalancerARN *string `json:"loadBalancerARN,omitempty"`

	Port *int64 `json:"port,omitempty"`

	Protocol *string `json:"protocol,omitempty"`

	SSLPolicy *string `json:"sslPolicy,omitempty"`
}

type LoadBalancerAddress struct {
	AllocationID *string `json:"allocationID,omitempty"`

	IPv6Address *string `json:"iPv6Address,omitempty"`

	IPAddress *string `json:"ipAddress,omitempty"`

	PrivateIPv4Address *string `json:"privateIPv4Address,omitempty"`
}

type LoadBalancerState struct {
	Code *string `json:"code,omitempty"`

	Reason *string `json:"reason,omitempty"`
}

type LoadBalancer_SDK struct {
	AvailabilityZones []*AvailabilityZone `json:"availabilityZones,omitempty"`

	CanonicalHostedZoneID *string `json:"canonicalHostedZoneID,omitempty"`

	CreatedTime *metav1.Time `json:"createdTime,omitempty"`

	CustomerOwnedIPv4Pool *string `json:"customerOwnedIPv4Pool,omitempty"`

	DNSName *string `json:"dnsName,omitempty"`

	IPAddressType *string `json:"ipAddressType,omitempty"`

	LoadBalancerARN *string `json:"loadBalancerARN,omitempty"`

	LoadBalancerName *string `json:"loadBalancerName,omitempty"`

	Scheme *string `json:"scheme,omitempty"`

	SecurityGroups []*string `json:"securityGroups,omitempty"`
	// Information about the state of the load balancer.
	State *LoadBalancerState `json:"state,omitempty"`

	Type *string `json:"type_,omitempty"`

	VPCID *string `json:"vpcID,omitempty"`
}

type Matcher struct {
	GrpcCode *string `json:"grpcCode,omitempty"`

	HTTPCode *string `json:"httpCode,omitempty"`
}

type PathPatternConditionConfig struct {
	Values []*string `json:"values,omitempty"`
}

type QueryStringConditionConfig struct {
	Values []*QueryStringKeyValuePair `json:"values,omitempty"`
}

type QueryStringKeyValuePair struct {
	Key *string `json:"key,omitempty"`

	Value *string `json:"value,omitempty"`
}

type RedirectActionConfig struct {
	Host *string `json:"host,omitempty"`

	Path *string `json:"path,omitempty"`

	Port *string `json:"port,omitempty"`

	Protocol *string `json:"protocol,omitempty"`

	Query *string `json:"query,omitempty"`

	StatusCode *string `json:"statusCode,omitempty"`
}

type RuleCondition struct {
	Field *string `json:"field,omitempty"`
	// Information about a host header condition.
	HostHeaderConfig *HostHeaderConditionConfig `json:"hostHeaderConfig,omitempty"`
	// Information about an HTTP header condition.
	//
	// There is a set of standard HTTP header fields. You can also define custom
	// HTTP header fields.
	HTTPHeaderConfig *HTTPHeaderConditionConfig `json:"httpHeaderConfig,omitempty"`
	// Information about an HTTP method condition.
	//
	// HTTP defines a set of request methods, also referred to as HTTP verbs. For
	// more information, see the HTTP Method Registry (https://www.iana.org/assignments/http-methods/http-methods.xhtml).
	// You can also define custom HTTP methods.
	HTTPRequestMethodConfig *HTTPRequestMethodConditionConfig `json:"httpRequestMethodConfig,omitempty"`
	// Information about a path pattern condition.
	PathPatternConfig *PathPatternConditionConfig `json:"pathPatternConfig,omitempty"`
	// Information about a query string condition.
	//
	// The query string component of a URI starts after the first '?' character
	// and is terminated by either a '#' character or the end of the URI. A typical
	// query string contains key/value pairs separated by '&' characters. The allowed
	// characters are specified by RFC 3986. Any character can be percentage encoded.
	QueryStringConfig *QueryStringConditionConfig `json:"queryStringConfig,omitempty"`
	// Information about a source IP condition.
	//
	// You can use this condition to route based on the IP address of the source
	// that connects to the load balancer. If a client is behind a proxy, this is
	// the IP address of the proxy not the IP address of the client.
	SourceIPConfig *SourceIPConditionConfig `json:"sourceIPConfig,omitempty"`

	Values []*string `json:"values,omitempty"`
}

type RulePriorityPair struct {
	Priority *int64 `json:"priority,omitempty"`

	RuleARN *string `json:"ruleARN,omitempty"`
}

type Rule_SDK struct {
	Actions []*Action `json:"actions,omitempty"`

	Conditions []*RuleCondition `json:"conditions,omitempty"`

	IsDefault *bool `json:"isDefault,omitempty"`

	Priority *string `json:"priority,omitempty"`

	RuleARN *string `json:"ruleARN,omitempty"`
}

type SSLPolicy struct {
	Name *string `json:"name,omitempty"`
}

type SourceIPConditionConfig struct {
	Values []*string `json:"values,omitempty"`
}

type SubnetMapping struct {
	AllocationID *string `json:"allocationID,omitempty"`

	IPv6Address *string `json:"iPv6Address,omitempty"`

	PrivateIPv4Address *string `json:"privateIPv4Address,omitempty"`

	SubnetID *string `json:"subnetID,omitempty"`
}

type Tag struct {
	Key *string `json:"key,omitempty"`

	Value *string `json:"value,omitempty"`
}

type TagDescription struct {
	Tags []*Tag `json:"tags,omitempty"`
}

type TargetDescription struct {
	AvailabilityZone *string `json:"availabilityZone,omitempty"`

	Port *int64 `json:"port,omitempty"`
}

type TargetGroupStickinessConfig struct {
	DurationSeconds *int64 `json:"durationSeconds,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`
}

type TargetGroupTuple struct {
	TargetGroupARN *string `json:"targetGroupARN,omitempty"`

	Weight *int64 `json:"weight,omitempty"`
}

type TargetGroup_SDK struct {
	HealthCheckEnabled *bool `json:"healthCheckEnabled,omitempty"`

	HealthCheckIntervalSeconds *int64 `json:"healthCheckIntervalSeconds,omitempty"`

	HealthCheckPath *string `json:"healthCheckPath,omitempty"`

	HealthCheckPort *string `json:"healthCheckPort,omitempty"`

	HealthCheckProtocol *string `json:"healthCheckProtocol,omitempty"`

	HealthCheckTimeoutSeconds *int64 `json:"healthCheckTimeoutSeconds,omitempty"`

	HealthyThresholdCount *int64 `json:"healthyThresholdCount,omitempty"`

	LoadBalancerARNs []*string `json:"loadBalancerARNs,omitempty"`
	// The codes to use when checking for a successful response from a target. If
	// the protocol version is gRPC, these are gRPC codes. Otherwise, these are
	// HTTP codes.
	Matcher *Matcher `json:"matcher,omitempty"`

	Port *int64 `json:"port,omitempty"`

	Protocol *string `json:"protocol,omitempty"`

	ProtocolVersion *string `json:"protocolVersion,omitempty"`

	TargetGroupARN *string `json:"targetGroupARN,omitempty"`

	TargetGroupName *string `json:"targetGroupName,omitempty"`

	TargetType *string `json:"targetType,omitempty"`

	UnhealthyThresholdCount *int64 `json:"unhealthyThresholdCount,omitempty"`

	VPCID *string `json:"vpcID,omitempty"`
}

type TargetHealthDescription struct {
	HealthCheckPort *string `json:"healthCheckPort,omitempty"`
}
